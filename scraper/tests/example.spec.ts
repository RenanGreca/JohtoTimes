import { test, expect } from '@playwright/test';
import { JohtoTimesPOM } from '../pom/JohtoTimes';
import { slugify } from '../pom/Slug';
import { Post } from '../pom/Post';
import { readFileSync, writeFileSync } from 'fs';

const vols = [2];
const issues = [21, 30];

for (const vol of vols) {
  for (const issue of issues) {
    test.describe(`JohtoTimes vol ${vol}, issue ${issue}`, () => {
      let pom: JohtoTimesPOM;
      let json: Post;
      test.beforeAll(async ({ browser }) => {
        const page = await browser.newPage();
        pom = new JohtoTimesPOM(page, vol, issue);
        await pom.goTo();

        json = {}
        json.volume = vol;
        json.issue = issue;

        writeFileSync(`./jsons/${vol}-${issue}.json`, JSON.stringify(json));
      });

      // Read and write JSON before and after each test
      test.beforeEach(() => {
        json = JSON.parse(readFileSync(`./jsons/${vol}-${issue}.json`, 'utf8'));
      })
      test.afterEach(() => {
        writeFileSync(`./jsons/${vol}-${issue}.json`, JSON.stringify(json));
      })

      test(`Get title and date`, async() => {
        await pom.preprocess();
        const title = await pom.getTitle();
        const date = await pom.getDate();

        json.title = title;
        json.date = date;
        json.slug = `${date}-${slugify(title)}`;

        const intro = await pom.getIntro();
        json.intro = intro;

        const news = await pom.getNews();
        json.news = news;

        const body = await pom.getBody();
        json.body = body;

        const mailbag = await pom.getMailbag();
        json.mailbag = mailbag;
      })

    })
  }
}
