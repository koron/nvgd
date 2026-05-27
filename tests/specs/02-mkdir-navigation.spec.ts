import { expect, test } from '@playwright/test';
import { OpfsPage } from '../pages/OpfsPage';
import { captureDialog, recordDialogs } from '../helpers/dialogs';
import { clearOPFS, existsOPFS, seedDirectory } from '../helpers/opfs';

test.describe('B. mkdir & directory navigation', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/opfs/');
    await clearOPFS(page);
  });

  test('B1: creating a directory adds a row and clears the input', async ({
    page,
  }) => {
    const opfs = new OpfsPage(page);
    await opfs.goto();

    await opfs.mkdir('foo');

    await expect(opfs.row('foo/')).toBeVisible();
    await expect(opfs.row('foo/')).toContainText('dir');
    await expect(opfs.mkdirName).toHaveValue('');
    expect(await existsOPFS(page, 'foo')).toBe('dir');
  });

  test('B2: empty mkdir shows an alert and creates nothing', async ({
    page,
  }) => {
    const opfs = new OpfsPage(page);
    await opfs.goto();

    const dialogPromise = captureDialog(page, 'accept');
    await opfs.mkdir('');

    const dialog = await dialogPromise;
    expect(dialog.message()).toMatch(/Need directory name/i);
    expect(await opfs.rowCount()).toBe(0);
  });

  test('B3: creating a same-named directory does not duplicate', async ({
    page,
  }) => {
    const opfs = new OpfsPage(page);
    await opfs.goto();

    await opfs.mkdir('dup');
    // The second mkdir succeeds idempotently (getDirectoryHandle with
    // create:true on an existing dir returns the same handle).
    await opfs.mkdir('dup');

    // Still only a single row visible.
    const matchingRows = await opfs.grid
      .locator('.grid-row')
      .locator('input.selectedFile[name="dup/"]')
      .count();
    expect(matchingRows).toBe(1);
  });

  test('B4: directory listing is sorted in natural order', async ({ page }) => {
    const opfs = new OpfsPage(page);
    await opfs.goto();

    // Create out of order; expect natural-numeric order in the UI.
    for (const name of ['file10', 'file2', 'file1']) {
      await opfs.mkdir(name);
    }

    const names = await opfs.rowNames();
    expect(names).toEqual(['file1/', 'file2/', 'file10/']);
  });

  test('B5: clicking a directory cd-s into it and updates breadcrumb', async ({
    page,
  }) => {
    const opfs = new OpfsPage(page);
    await opfs.goto();

    await opfs.mkdir('childdir');
    await opfs.openDirectory('childdir');

    expect(await opfs.breadcrumbSegments()).toEqual(['(Root)', 'childdir']);
    await expect(page).toHaveTitle(/childdir/);
    expect(await opfs.rowCount()).toBe(0);
  });

  test('B6: clicking a parent crumb returns to that level', async ({
    page,
  }) => {
    const opfs = new OpfsPage(page);
    await seedDirectory(page, 'a/b/c');
    await opfs.goto('a/b/c/');

    expect(await opfs.breadcrumbSegments()).toEqual([
      '(Root)',
      'a',
      'b',
      'c',
    ]);

    // Click the "a" crumb (an <a>) to jump back two levels.
    await opfs.breadcrumb.getByText('a', { exact: true }).click();

    expect(await opfs.breadcrumbSegments()).toEqual(['(Root)', 'a']);
  });

  test('B7: browser back/forward replays the directory history', async ({
    page,
  }) => {
    const opfs = new OpfsPage(page);
    await seedDirectory(page, 'level1/level2');
    await opfs.goto();

    await opfs.openDirectory('level1');
    await opfs.openDirectory('level2');
    expect(await opfs.breadcrumbSegments()).toEqual([
      '(Root)',
      'level1',
      'level2',
    ]);

    await page.goBack();
    await expect.poll(() => opfs.breadcrumbSegments()).toEqual([
      '(Root)',
      'level1',
    ]);

    await page.goBack();
    await expect.poll(() => opfs.breadcrumbSegments()).toEqual(['(Root)']);

    await page.goForward();
    await expect.poll(() => opfs.breadcrumbSegments()).toEqual([
      '(Root)',
      'level1',
    ]);
  });

  test('B8: F5 reload preserves the current path via hash', async ({
    page,
  }) => {
    const opfs = new OpfsPage(page);
    await seedDirectory(page, 'deep/nested/dir');
    await opfs.goto('deep/nested/dir/');

    // Track that a *future* reload restoration works without alerts.
    const log = recordDialogs(page, 'accept');
    await page.reload();
    log.stop();
    expect(log.messages).toEqual([]);

    expect(await opfs.breadcrumbSegments()).toEqual([
      '(Root)',
      'deep',
      'nested',
      'dir',
    ]);
  });
});
