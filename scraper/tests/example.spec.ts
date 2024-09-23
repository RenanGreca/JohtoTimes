import { test } from '@playwright/test';
import { JohtoTimesPOM } from '../pom/JohtoTimes';
import { slugify } from '../pom/Slug';
import { issueMarkdown, mailbagMarkdown, newsMarkdown, Post, postMarkdown } from '../pom/Post';
import { readFileSync, writeFileSync } from 'fs';


// const vols = {
//   1: Array.from({ length: 48 }, (_, i) => i + 1),
//   2: Array.from({ length: 35 }, (_, i) => i + 1),
// }
//
const vols = {
  2: [36, 37, 38]
}

for (const v in vols) {
  const vol = Number(v);
  const issues = vols[v];
  for (const issue of issues) {
    test.describe(`JohtoTimes vol ${vol}, issue ${issue}`, () => {
      test(`Get data from Substack`, async ({ browser }) => {

        const pom = await test.step("Load page in Substack", async () => {
          const page = await browser.newPage();
          const pom = new JohtoTimesPOM(page, vol, issue);
          await pom.goTo();
          await pom.preprocess();
          return pom;
        })

        await test.step("Parse data from webpage", async () => {
          const json = {} as Post;
          json.volume = vol;
          json.issue = issue;

          json.title = await pom.getTitle();
          json.date = await pom.getDate();
          json.slug = `${json.date}-${slugify(json.title)}`;

          json.description = await pom.getDescription();
          json.intro = await pom.getIntro();
          json.news = await pom.getNews();
          json.body = await pom.getBody();
          json.img = await pom.getImg();
          json.mailbag = await pom.getMailbag();
          json.comments = await pom.getComments();

          writeFileSync(`./jsons/${vol}-${issue}.json`, JSON.stringify(json));
          return json;
        })

        const post = await test.step(`Parse JSON`, () => {
          console.log(`Parsing: ./jsons/${vol}-${issue}.json`);
          const post = JSON.parse(readFileSync(`./jsons/${vol}-${issue}.json`, 'utf8')) as Post;
          return post;
        })

        await test.step(`Generate Markdowns`, () => {
          const postMD = postMarkdown(post);
          writeFileSync(`./posts/${post.slug}.md`, postMD);

          const issueMD = issueMarkdown(post);
          writeFileSync(`./issues/${post.slug}.md`, issueMD);

          const newsMD = newsMarkdown(post);
          writeFileSync(`./news/${post.date}-news.md`, newsMD);

          const mailbagMD = mailbagMarkdown(post);
          writeFileSync(`./mailbag/${post.date}-mailbag.md`, mailbagMD);
        })
      })
    })
  }
}
