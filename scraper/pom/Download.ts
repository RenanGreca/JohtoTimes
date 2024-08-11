import { createWriteStream } from "fs";
import { slugify } from "./Slug";
import axios from "axios";
import { resolve } from "path";

export async function download(url: string, name: string) {
  const slug = slugify(name).slice(0, 100);
  const extension = url.split('.').pop();
  const filename = `${slug}.${extension}`;
  const path = resolve('downloads', filename);
   
  const response = await axios({
    url,
    method: 'GET',
    responseType: 'stream',
  })
  const writer = createWriteStream(path);
  response.data.pipe(writer);

  return new Promise<string>((resolve, reject) => {
    writer.on('finish', () => {
      resolve(filename);
    });
    writer.on('error', reject);
  })
}
