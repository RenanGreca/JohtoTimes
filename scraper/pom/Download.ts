import { createWriteStream } from "fs";
import { slugify } from "./Slug";
import axios from "axios";
import { resolve } from "path";

export async function download(url: string, name: string) {
  const slug = slugify(name).slice(0, 100);
  const extension = url.split('.').pop();
  const path = resolve('downloads', `${slug}.${extension}`);
   
  const response = await axios({
    url,
    method: 'GET',
    responseType: 'stream',
  })
  const writer = createWriteStream(path);
  response.data.pipe(writer);

  return new Promise((resolve, reject) => {
    writer.on('finish', resolve);
    writer.on('error', reject);
  })
}
