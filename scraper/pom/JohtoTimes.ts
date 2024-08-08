import { Page, Locator } from "@playwright/test";
import TurndownService from 'turndown';
import { download } from "./Download";
import { assert } from "console";

export class JohtoTimesPOM {
  public static readonly POM_NAME = "JohtoTimes";

  private url: string;
  private components: string[];
  private body: Locator;
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
    const title = await locator.innerText();

    if (title.includes(' - ')) {
      return title.split(' - ')[1];
    }
    return title;
  }

  async getDate() {
    const locator = this.page.locator('.post-header').locator('.pc-reset').nth(1)
    const datestring = await locator.innerText();
    const date = new Date(`${datestring} 00:00:00 UTC`);

    const year = date.getFullYear().toString();
    let month = (date.getMonth() + 1).toString();
    if (date.getMonth() < 9) {
      month = `0${month}`;
    }
    let day = date.getDate().toString();
    if (date.getDate() < 10) {
      day = `0${day}`;
    }
    return `${year}-${month}-${day}`;
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
        const img = await this.downloadImage(arr, i);
        ret.push(img);
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
    const imgstart = arr[i]
    const imgblock = arr[i + 2];
    const imgend = arr[i + 6];
    const imgtitle = arr[i + 8];
    const img = `${imgstart}${imgblock}${imgend}*${imgtitle}*`;
    const imgurl = imgend.split('](')[1].replace(')', '');

    await download(imgurl, imgtitle);
    return img;
  }

}
