import { expect, test } from '@playwright/test';
import { OpfsPage } from '../pages/OpfsPage';
import { captureDialog } from '../helpers/dialogs';
import { clearOPFS, seedDirectory } from '../helpers/opfs';

test.describe('A. Initial load & title', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/opfs/');
    await clearOPFS(page);
  });

  test('A1: opens /opfs/ with empty root', async ({ page }) => {
    const opfs = new OpfsPage(page);
    await opfs.goto();

    await expect(page).toHaveTitle(/^OPFS:\s+\(Root\)\/?$/);
    await expect(opfs.breadcrumb).toContainText('(Root)');

    // grid header is always present; there should be no data rows.
    await expect(opfs.grid.locator('.grid-header')).toBeVisible();
    expect(await opfs.rowCount()).toBe(0);

    // Action buttons that require a selection start disabled.
    await expect(opfs.deleteBtn).toBeDisabled();
    await expect(opfs.duckdbBtn).toBeDisabled();
  });

  test('A2: opens with hash path after seeding sub1/sub2', async ({ page }) => {
    const opfs = new OpfsPage(page);
    // Seed the directory tree from inside the page origin.
    await seedDirectory(page, 'sub1/sub2');

    await opfs.goto('sub1/sub2/');

    await expect(page).toHaveTitle(/sub1\/sub2/);

    const segments = await opfs.breadcrumbSegments();
    expect(segments).toEqual(['(Root)', 'sub1', 'sub2']);
  });

  test('A3: invalid hash path raises an alert and falls back to Root', async ({
    page,
  }) => {
    const opfs = new OpfsPage(page);
    // No seed: navigating to /opfs/#missing/ will fail getDirectoryHandle.
    const dialogPromise = captureDialog(page, 'accept');
    await page.goto('/opfs/#missing/');

    const dialog = await dialogPromise;
    expect(dialog.message().toLowerCase()).toMatch(/notfound|error/);

    // Even after the alert, the grid header still renders.
    await expect(opfs.grid.locator('.grid-header')).toBeVisible();
  });
});
