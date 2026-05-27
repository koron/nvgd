import { expect, test } from '@playwright/test';
import { OpfsPage } from '../pages/OpfsPage';
import { recordDialogs } from '../helpers/dialogs';
import { clearOPFS, seedDirectory } from '../helpers/opfs';

test.describe('A. Initial load & title', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/opfs/');
    await clearOPFS(page);
  });

  test('A1: opens /opfs/ with empty root', async ({ page }) => {
    const opfs = new OpfsPage(page);
    await opfs.goto();

    // The root FileSystemDirectoryHandle has an empty name, so the
    // title path collapses to "/". The breadcrumb separately uses the
    // friendly label "(Root)" only inside #header.
    await expect(page).toHaveTitle(/^OPFS:\s+\/$/);
    await expect(opfs.breadcrumb).toContainText('(Root)');

    await expect(opfs.grid.locator('.grid-header')).toBeVisible();
    expect(await opfs.rowCount()).toBe(0);

    await expect(opfs.deleteBtn).toBeDisabled();
    await expect(opfs.duckdbBtn).toBeDisabled();
  });

  test('A2: opens with hash path after seeding sub1/sub2', async ({ page }) => {
    const opfs = new OpfsPage(page);
    await seedDirectory(page, 'sub1/sub2');

    await opfs.goto('sub1/sub2/');

    await expect(page).toHaveTitle(/sub1\/sub2/);

    const segments = await opfs.breadcrumbSegments();
    expect(segments).toEqual(['(Root)', 'sub1', 'sub2']);
  });

  test('A3: invalid hash path stays silent (no alert, init bails out)', async ({
    page,
  }) => {
    // Document the *actual* behaviour: init() awaits setCurrPath() and
    // that call throws NotFoundError for a missing directory. The
    // promise is never caught, so the app emits no alert. This test
    // guards against a future regression that adds a stray alert.
    const log = recordDialogs(page, 'accept');
    await page.goto('/opfs/#__no_such_dir__/');
    await page.waitForTimeout(500);
    const messages = log.stop();
    expect(messages).toEqual([]);
  });
});
