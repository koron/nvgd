import { Page } from '@playwright/test';

/** OPFS UI のルートに移動する */
export async function gotoOPFS(page: Page): Promise<void> {
  await page.goto('/opfs/');
}

/**
 * OPFS にファイルを直接作成する（page.evaluate 経由）。
 * gotoOPFS() で対象オリジンに移動済みであること。
 */
export async function createOPFSFile(
  page: Page,
  name: string,
  content: string,
): Promise<void> {
  await page.evaluate(
    async ({ name, content }) => {
      const root = await navigator.storage.getDirectory();
      const fh = await root.getFileHandle(name, { create: true });
      const w = await fh.createWritable();
      await w.write(content);
      await w.close();
    },
    { name, content },
  );
}

/**
 * OPFS にサブディレクトリを直接作成する（page.evaluate 経由）。
 * gotoOPFS() で対象オリジンに移動済みであること。
 */
export async function createOPFSDir(page: Page, name: string): Promise<void> {
  await page.evaluate(async (name) => {
    const root = await navigator.storage.getDirectory();
    await root.getDirectoryHandle(name, { create: true });
  }, name);
}

/**
 * OPFS のサブディレクトリ内にファイルを作成する。
 * gotoOPFS() で対象オリジンに移動済みであること。
 */
export async function createOPFSFileInDir(
  page: Page,
  dir: string,
  name: string,
  content: string,
): Promise<void> {
  await page.evaluate(
    async ({ dir, name, content }) => {
      const root = await navigator.storage.getDirectory();
      const dh = await root.getDirectoryHandle(dir, { create: true });
      const fh = await dh.getFileHandle(name, { create: true });
      const w = await fh.createWritable();
      await w.write(content);
      await w.close();
    },
    { dir, name, content },
  );
}

/** Reload ボタンをクリックしてファイル一覧を更新する */
export async function reloadListing(page: Page): Promise<void> {
  await page.click('#command-reload');
}

/** OPFS からファイルの内容を読み取る（検証用） */
export async function readOPFSFile(page: Page, name: string): Promise<string> {
  return page.evaluate(async (name) => {
    const root = await navigator.storage.getDirectory();
    const fh = await root.getFileHandle(name);
    const file = await fh.getFile();
    return file.text();
  }, name);
}
