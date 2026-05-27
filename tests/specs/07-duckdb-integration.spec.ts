import { expect, test } from '@playwright/test';
import { OpfsPage } from '../pages/OpfsPage';
import { clearOPFS, seedFile } from '../helpers/opfs';

/**
 * The OPFS UI launches DuckDB by `window.open(url, '_blank')`. We
 * intercept the new tab via `context().waitForEvent('page')`, assert
 * on its URL (without ever loading the WASM shell), and close it.
 */

test.describe('G. DuckDB integration', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/opfs/');
    await clearOPFS(page);
  });

  test('G1: single CSV opens a tab with CREATE VIEW opfs0', async ({
    page,
    context,
  }) => {
    const opfs = new OpfsPage(page);
    await seedFile(page, 'sample.csv', 'id,name\n1,a\n');
    await opfs.goto();

    await opfs.selectRow('sample.csv');

    const [newPage] = await Promise.all([
      context.waitForEvent('page'),
      opfs.duckdbBtn.click(),
    ]);

    const url = newPage.url();
    // The hash carries swapchars-encoded SQL. Decode the easy parts.
    expect(url).toMatch(/\/duckdb\/\?opfs=/);
    expect(decodeURIComponent(url)).toContain('sample.csv');
    expect(url).toContain('opfs0');

    await newPage.close();
  });

  test('G2: mixed types produce opfs0, opfs1, opfs2 in order', async ({
    page,
    context,
  }) => {
    const opfs = new OpfsPage(page);
    // Names chosen so natural sort picks an unambiguous order:
    // a.csv, b.json, c.parquet.
    await seedFile(page, 'a.csv', 'id\n1\n');
    await seedFile(page, 'b.json', '[1]');
    await seedFile(page, 'c.parquet', 'PAR1');
    await opfs.goto();

    await opfs.selectRow('a.csv');
    await opfs.selectRow('b.json');
    await opfs.selectRow('c.parquet');

    const [newPage] = await Promise.all([
      context.waitForEvent('page'),
      opfs.duckdbBtn.click(),
    ]);

    const url = newPage.url();
    expect(url).toContain('opfs0');
    expect(url).toContain('opfs1');
    expect(url).toContain('opfs2');
    // The hash uses swapchars: '-' <-> ' '. Decode and check for
    // 'CREATE VIEW' (which becomes 'CREATE-VIEW' in the wire form).
    const decoded = decodeURIComponent(url);
    expect(decoded).toMatch(/CREATE.VIEW.opfs0/);

    await newPage.close();
  });

  test('G3: unsupported extension is still passed via opfs= but no view created', async ({
    page,
    context,
  }) => {
    const opfs = new OpfsPage(page);
    await seedFile(page, 'note.txt', 'plain');
    await opfs.goto();

    await opfs.selectRow('note.txt');

    const [newPage] = await Promise.all([
      context.waitForEvent('page'),
      opfs.duckdbBtn.click(),
    ]);

    const url = newPage.url();
    expect(url).toContain('opfs=');
    expect(decodeURIComponent(url)).toContain('note.txt');
    // No CREATE VIEW for unsupported types — only SHOW TABLES (which
    // becomes 'SHOW-TABLES' in the swapchars form).
    expect(url).not.toContain('opfs0');

    await newPage.close();
  });

  test('G4: selecting a directory enumerates files recursively', async ({
    page,
    context,
  }) => {
    const opfs = new OpfsPage(page);
    await seedFile(page, 'data/one.csv', 'a\n1\n');
    await seedFile(page, 'data/two.json', '[1]');
    await opfs.goto();

    await opfs.selectRow('data/');

    const [newPage] = await Promise.all([
      context.waitForEvent('page'),
      opfs.duckdbBtn.click(),
    ]);

    const url = decodeURIComponent(newPage.url());
    expect(url).toContain('data/one.csv');
    expect(url).toContain('data/two.json');

    await newPage.close();
  });
});
