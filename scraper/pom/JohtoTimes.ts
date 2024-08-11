import { Page, Locator } from "@playwright/test";
import TurndownService from 'turndown';
import { download } from "./Download";
import { assert } from "console";
import { Comment } from "./Post";
import { dateToString, escapeString } from "./Slug";

export class JohtoTimesPOM {
  public static readonly POM_NAME = "JohtoTimes";

  private url: string;
  private components: string[];
  private body: Locator;

  // The first (header) image in the post body
  firstImage: string;
  constructor(private page: Page, vol: number, issue: number) {
    this.url = `https://johto.substack.com/p/vol${vol}-${issue}`;
    this.body = this.page.locator('.body');
  }

  async goTo() {
    await this.page.goto(this.url);
  }

  async preprocess() {
    const html = await this.body.innerHTML();
    this.components = html.split("<div><hr></div>");
  }

  async getTitle() {
    const locator = this.page.locator('.post-title')
    let title = await locator.innerText();

    if (title.includes(' - ')) {
      title = title.split(' - ')[1];
    }

    return escapeString(title);
  }

  async getDate() {
    const locator = this.page.locator('.post-header').locator('.pc-reset').nth(1)
    const datestring = await locator.innerText();
    const date = new Date(`${datestring} 00:00:00 UTC`);
    return dateToString(date);
  }

  async getDescription() {
    const locator = this.page.locator('.subtitle')
    const subtitle = await locator.innerText();

    return escapeString(subtitle);
  }

  async getIntro() {
    const turndownService = new TurndownService();
    const md = turndownService.turndown(this.components[0]);

    return md.split('\n').filter((el: string) => el.length > 0);
  }

  async getNews() {
    const turndownService = new TurndownService();
    const md = turndownService.turndown(this.components[1]);
    return this.filterMarkdown(md, "News");
  }

  async getBody() {
    const mailbagIndex = this.sectionIndex("Mailbag");
    const featureIndex = this.sectionIndex("Feature");
    assert(featureIndex > -1);
    let body = this.components.slice(featureIndex, this.components.length-2).join("<div><hr></div>");
    if (mailbagIndex > -1) {
      body = this.components.slice(featureIndex, mailbagIndex).join("<hr>");
    }

    const turndownService = new TurndownService();
    const md = turndownService.turndown(body);
    return this.filterMarkdown(md, "Feature");
  }

  async getImg() {
    const mailbagIndex = this.sectionIndex("Mailbag");
    const featureIndex = this.sectionIndex("Feature");
    assert(featureIndex > -1);
    let body = this.components.slice(featureIndex, this.components.length-2).join("<div><hr></div>");
    if (mailbagIndex > -1) {
      body = this.components.slice(featureIndex, mailbagIndex).join("<hr>");
    }

    const turndownService = new TurndownService();
    const md = turndownService.turndown(body);
  
    let arr = md.split('\n')
    for (let i = 0; i < arr.length; i++) {
      if (this.isImage(arr, i)) {
        const [_, imgpath] = await this.downloadImage(arr, i);
        console.log("Found image", imgpath);
        return imgpath;
      }
    }

    return "";
  }

  async getMailbag() {
    const mailbagIndex = this.sectionIndex("Mailbag");
    if (mailbagIndex === -1) {
      return [];
    }
    const body = this.components.slice(mailbagIndex, this.components.length-1).join("<div><hr></div>");

    const turndownService = new TurndownService();
    const md = turndownService.turndown(body);
    return this.filterMarkdown(md, "Mailbag");
  }

  async getComments() {
    const moreComments = this.page.locator('.more-comments');
    if (await moreComments.isVisible()) {
      await moreComments.click();
    }

    const commentsBlock = this.page.locator('.comment-rest');

    let comments: Comment[] = [];

    for (const commentBlock of await commentsBlock.all()) {
      const comment = {} as Comment;
      comment.name = await commentBlock.locator('.commenter-name').innerText();

      const body = await commentBlock.locator('.comment-body').innerText();
      const turndownService = new TurndownService();
      const md = turndownService.turndown(body);
      comment.body = md.split('\n').filter((el: string) => el.length > 0);

      const datestring = await commentBlock.locator('.comment-timestamp').innerText();
      const date = new Date(`${datestring} 00:00:00 UTC`);
      comment.date = dateToString(date);

      comments.push(comment);
    }

    return comments;
  }

  /**
   * Gets the index of a section in the HTML blocks.
   * @param name The name of the section
   * @returns The index of the section, or -1 if not found
   */
  private sectionIndex(name: string) {
    for (let i = 0; i < this.components.length; i++) {
      const el = this.components[i];
      if (el.startsWith(`<h4 class="header-anchor-post">${name}`)) {
        return i;
      }
    }

    return -1;
  }

  /**
   * Filter the markdown to remove empty blocks, the title block,
   * and organize image blocks.
   * @param md The markdown to filter
   * @param section The section to filter
   * @returns The filtered markdown
   */
  private async filterMarkdown(md: string, section: string) {
    let arr = md.split('\n')
    let ret: string[] = [];
    for (let i = 0; i < arr.length; i++) {
      const el = arr[i];
      if (el.length == 0) {
        continue;
      } else
      if (el.startsWith(`#### ${section}`)) {
        continue;
      } else
      if (this.isImage(arr, i)) {
        const [imgmd, _] = await this.downloadImage(arr, i);
        ret.push("");
        ret.push(imgmd);
        ret.push("");
        i += 8;
      } else {
        ret.push(el);
      }
    }

    return ret;
  }

  /**
   * Check if the element at index i is an image
   * @param arr The array of markdown blocks
   * @param i The index to check
   * @returns True if the index contains an image, false otherwise
   */
  private isImage(arr: string[], i: number) {
    if (arr[i] === "[") {
      assert(i + 2 < arr.length, `i + 2 is out of bounds`);
      return arr[i] === "[" && arr[i + 2].startsWith("![");
    }
    return false;
  }

  /**
   * Given an index to the array of markdown blocks,
   * "clean up" the image markdown and download a copy of the image.
   * @param arr The array of markdown blocks
   * @param i The index of the image
   * @returns The markdown for the image
   */
  private async downloadImage(arr: string[], i: number) {
    assert(i + 8 < arr.length, "i + 8 is out of bounds");
    // const imgstart = arr[i]
    // const imgblock = arr[i + 2];
    const imgend = arr[i + 6];
    const imgtitle = arr[i + 8];
    const imgurl = imgend.split('](')[1].replace(')', '');

    const imgpath = "/web/images/" + await download(imgurl, imgtitle);

    if (!this.firstImage) {
      this.firstImage = imgpath;
    }

    const imgmd = `[![${imgtitle}](${imgpath})](${imgpath})*${imgtitle}*`;

    return [imgmd, imgpath];
  }

}
