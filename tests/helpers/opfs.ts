import type { Page } from '@playwright/test';

/**
 * Wipe every entry under the OPFS root. Safe to call before each test
 * as a defensive measure even though Playwright's isolated contexts
 * already give us a fresh storage. Must be invoked on a page already
 * navigated to the nvgd origin.
 */
export async function clearOPFS(page: Page): Promise<void> {
  await page.evaluate(async () => {
    const root = await navigator.storage.getDirectory();
    // `entries()` is async-iterable but TypeScript's lib.dom.d.ts
    // sometimes lags; cast through `any`.
    for await (const [name] of (root as any).entries()) {
      await root.removeEntry(name, { recursive: true });
    }
  });
}

/**
 * Create a directory tree inside OPFS, returning the deepest handle's
 * path. Intermediate directories are created as needed.
 */
export async function seedDirectory(page: Page, path: string): Promise<void> {
  await page.evaluate(async (p: string) => {
    const segments = p.replace(/^\/+|\/+$/g, '').split('/').filter(Boolean);
    let dir = await navigator.storage.getDirectory();
    for (const seg of segments) {
      dir = await dir.getDirectoryHandle(seg, { create: true });
    }
  }, path);
}

/**
 * Write a UTF-8 file into OPFS at the given path. Intermediate
 * directories are created as needed.
 */
export async function seedFile(
  page: Page,
  path: string,
  content: string,
): Promise<void> {
  await page.evaluate(
    async ({ p, body }: { p: string; body: string }) => {
      const segments = p.replace(/^\/+/, '').split('/');
      const fileName = segments.pop()!;
      let dir = await navigator.storage.getDirectory();
      for (const seg of segments.filter(Boolean)) {
        dir = await dir.getDirectoryHandle(seg, { create: true });
      }
      const fh = await dir.getFileHandle(fileName, { create: true });
      const w = await fh.createWritable();
      await w.write(body);
      await w.close();
    },
    { p: path, body: content },
  );
}

/**
 * Write a binary file of a specific length into OPFS (used for the
 * 64KiB boundary test).
 */
export async function seedBinaryFile(
  page: Page,
  path: string,
  size: number,
): Promise<void> {
  await page.evaluate(
    async ({ p, n }: { p: string; n: number }) => {
      const segments = p.replace(/^\/+/, '').split('/');
      const fileName = segments.pop()!;
      let dir = await navigator.storage.getDirectory();
      for (const seg of segments.filter(Boolean)) {
        dir = await dir.getDirectoryHandle(seg, { create: true });
      }
      const fh = await dir.getFileHandle(fileName, { create: true });
      const w = await fh.createWritable();
      await w.write(new Uint8Array(n));
      await w.close();
    },
    { p: path, n: size },
  );
}

/**
 * List immediate children of the given OPFS directory (default: root).
 * Returns names sorted by browser default (mirrors what the UI shows).
 */
export async function listOPFS(
  page: Page,
  path: string = '/',
): Promise<string[]> {
  return await page.evaluate(async (p: string) => {
    const segments = p.replace(/^\/+|\/+$/g, '').split('/').filter(Boolean);
    let dir = await navigator.storage.getDirectory();
    for (const seg of segments) {
      dir = await dir.getDirectoryHandle(seg);
    }
    const names: string[] = [];
    for await (const [name] of (dir as any).entries()) {
      names.push(name);
    }
    names.sort((a, b) =>
      a.localeCompare(b, undefined, { numeric: true, sensitivity: 'base' }),
    );
    return names;
  }, path);
}

/**
 * Read a file's text contents from OPFS.
 */
export async function readOPFSFile(
  page: Page,
  path: string,
): Promise<string> {
  return await page.evaluate(async (p: string) => {
    const segments = p.replace(/^\/+/, '').split('/');
    const fileName = segments.pop()!;
    let dir = await navigator.storage.getDirectory();
    for (const seg of segments.filter(Boolean)) {
      dir = await dir.getDirectoryHandle(seg);
    }
    const fh = await dir.getFileHandle(fileName);
    const f = await fh.getFile();
    return await f.text();
  }, path);
}

/**
 * Read a file's size from OPFS.
 */
export async function statOPFSFile(
  page: Page,
  path: string,
): Promise<{ size: number }> {
  return await page.evaluate(async (p: string) => {
    const segments = p.replace(/^\/+/, '').split('/');
    const fileName = segments.pop()!;
    let dir = await navigator.storage.getDirectory();
    for (const seg of segments.filter(Boolean)) {
      dir = await dir.getDirectoryHandle(seg);
    }
    const fh = await dir.getFileHandle(fileName);
    const f = await fh.getFile();
    return { size: f.size };
  }, path);
}

/**
 * Check whether a path exists in OPFS, returning 'file', 'dir', or null.
 */
export async function existsOPFS(
  page: Page,
  path: string,
): Promise<'file' | 'dir' | null> {
  return await page.evaluate(async (p: string) => {
    const segments = p.replace(/^\/+|\/+$/g, '').split('/');
    const last = segments.pop()!;
    let dir = await navigator.storage.getDirectory();
    for (const seg of segments.filter(Boolean)) {
      try {
        dir = await dir.getDirectoryHandle(seg);
      } catch {
        return null;
      }
    }
    try {
      await dir.getDirectoryHandle(last);
      return 'dir' as const;
    } catch { /* not a dir */ }
    try {
      await dir.getFileHandle(last);
      return 'file' as const;
    } catch { /* not a file */ }
    return null;
  }, path);
}
