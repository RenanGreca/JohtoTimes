import { test } from '@playwright/test';
import { JohtoTimesPOM } from '../pom/JohtoTimes';
import { slugify } from '../pom/Slug';
import { issueMarkdown, mailbagMarkdown, newsMarkdown, Post, postMarkdown } from '../pom/Post';
import { readFileSync, writeFileSync } from 'fs';

const vols = [2];
const issues = [21, 29, 30];

for (const vol of vols) {
  for (const issue of issues) {
    test.describe(`JohtoTimes vol ${vol}, issue ${issue}`, () => {
      let pom: JohtoTimesPOM;
      let json: Post;
      test(`Get data from Substack`, async ({ browser }) => {
        const page = await browser.newPage();
        pom = new JohtoTimesPOM(page, vol, issue);
        await pom.goTo();

        json = {} as Post;
        json.volume = vol;
        json.issue = issue;

        await pom.preprocess();
        const title = await pom.getTitle();
        const date = await pom.getDate();

        json.title = title;
        json.date = date;
        json.slug = `${date}-${slugify(title)}`;

        const description = await pom.getDescription();
        json.description = description;

        const intro = await pom.getIntro();
        json.intro = intro;

        const news = await pom.getNews();
        json.news = news;

        const body = await pom.getBody();
        json.body = body;

        const img = await pom.getImg();
        json.img = img;

        const mailbag = await pom.getMailbag();
        json.mailbag = mailbag;

        const comments = await pom.getComments();
        json.comments = comments;

        writeFileSync(`./jsons/${vol}-${issue}.json`, JSON.stringify(json));
      })

      test(`Generate Markdowns from JSON`, () => {
        console.log(`Parsing: ./jsons/${vol}-${issue}.json`);
        const post = JSON.parse(readFileSync(`./jsons/${vol}-${issue}.json`, 'utf8')) as Post;

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
  }
}
