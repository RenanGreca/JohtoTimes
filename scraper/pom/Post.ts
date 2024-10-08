export interface Post {
  title?: string;
  date?: string;
  slug?: string;
  category?: string;
  img?: string;
  volume?: number;
  issue?: number;
  description?: string;
  intro?: string[];
  news?: string[];
  body?: string[];
  mailbag?: string[];
  comments?: Comment[];

}

export function postMarkdown(post: Post): string {
  let markdown = `---\n`;
  markdown += `Title: '${post.title}'\n`;
  markdown += `Volume: ${post.volume}\n`;
  markdown += `Issue: ${post.issue}\n`;
  markdown += `Header: '${post.img}'\n`;
  markdown += `Description: '${post.description}'\n`;
  markdown += `Slug: '${post.slug}'\n`;
  markdown += `Date: '${post.date}'\n`;
  markdown += `---\n`;

  markdown += `### Feature: ${post.title}\n`;

  markdown += post.body.join('\n\n');

  return markdown;
}

export function issueMarkdown(post: Post): string {
  let markdown = `---\n`;
  markdown += `Title: '${post.title}'\n`;
  // markdown += `News: '${post.date}-news'\n`;
  // markdown += `Post: '${post.date}-${post.slug}'\n`;
  // markdown += `Mailbag: '${post.date}-mailbag'\n`;
  markdown += `Volume: ${post.volume}\n`;
  markdown += `Issue: ${post.issue}\n`;
  markdown += `Header: '${post.img}'\n`;
  markdown += `Description: '${post.description}'\n`;
  markdown += `Slug: '${post.slug}'\n`;
  markdown += `Date: '${post.date}'\n`;
  markdown += `---\n`;

  markdown += post.intro.join('\n\n');

  return markdown;
}

export function newsMarkdown(post: Post): string {
  let markdown = `---\n`;
  markdown += `Title: ${post.date} News\n`;
  markdown += `Volume: ${post.volume}\n`;
  markdown += `Issue: ${post.issue}\n`;
  markdown += `Date: '${post.date}'\n`;
  markdown += `---\n`;

  markdown += post.news.join('\n\n');

  return markdown;
}

export function mailbagMarkdown(post: Post): string {
  let markdown = `---\n`;
  markdown += `Title: ${post.date} Mailbag\n`;
  markdown += `Volume: ${post.volume}\n`;
  markdown += `Issue: ${post.issue}\n`;
  markdown += `Date: '${post.date}'\n`;
  markdown += `---\n`;

  markdown += post.mailbag.join('\n\n');

  return markdown;
}

export interface Comment {
  name: string;
  date: string;
  body: string[];
}
