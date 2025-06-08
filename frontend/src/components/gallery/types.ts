interface Ithumb {
  sizeName: string;
  filePath: string;
}

export interface Ifile {
  id: number;
  originalName: string;
  mimeType: string;
  filePath: string;
  thumbnails: { [key: string]: Ithumb };
}
