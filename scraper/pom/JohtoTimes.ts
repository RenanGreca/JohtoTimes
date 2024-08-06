import { Page, Locator } from "@playwright/test";
import TurndownService from 'turndown';
import { download } from "./Download";

export class JohtoTimesPOM {
  public static readonly POM_NAME = "JohtoTimes";

  private url: string;
  private body: Locator;
  constructor(private page: Page, vol: number, issue: number) {
    this.url = `https://johto.substack.com/p/vol${vol}-${issue}`;
    this.body = this.page.locator('.body');
  }

  async goTo() {
    await this.page.goto(this.url);
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
    const html = await this.body.innerHTML();
    const turndownService = new TurndownService();
    const intro = html.split("<hr>")[0];
    const md = turndownService.turndown(intro);

    return md.split('\n').filter((el: string) => el.length > 0);
  }

  async getNews() {
    const html = await this.body.innerHTML();
    const turndownService = new TurndownService();
    const news = html.split("<hr>")[1];
    const md = turndownService.turndown(news);

    let arr = md.split('\n')
    let ret: string[] = [];
    for (let i = 0; i < arr.length; i++) {
      const el = arr[i];
      if (el.length == 0) {
        continue;
      } else
      if (el.startsWith('#### News')) {
        continue;
      } else
      if (el === "[" && arr[i + 2].startsWith("![")) {
        const img = `${el}${arr[i + 2]}${arr[i + 6]}*${arr[i + 8]}*`;
        const imgtitle = arr[i + 8];
        const imgurl = arr[i + 6].split('](')[1].replace(')', '');
        await download(imgurl, imgtitle);
        ret.push(img);
        i += 8;
      } else {
        ret.push(el);
      }
    }

    return ret;

  }

}
