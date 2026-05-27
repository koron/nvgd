import { expect, test } from '@playwright/test';
import { OpfsPage } from '../pages/OpfsPage';
import { clearOPFS, seedDirectory } from '../helpers/opfs';

test.describe('J. History & breadcrumb', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/opfs/');
    await clearOPFS(page);
  });

  test('J1: navigates three levels deep, then goes back to Root', async ({
    page,
  }) => {
    const opfs = new OpfsPage(page);
    await seedDirectory(page, 'l1/l2/l3');
    await opfs.goto();

    await opfs.openDirectory('l1');
    await opfs.openDirectory('l2');
    await opfs.openDirectory('l3');
    expect(await opfs.breadcrumbSegments()).toEqual([
      '(Root)',
      'l1',
      'l2',
      'l3',
    ]);

    await page.goBack();
    await page.goBack();
    await page.goBack();
    await expect.poll(() => opfs.breadcrumbSegments()).toEqual(['(Root)']);
  });

  test('J2: clicking a mid-level crumb jumps to that level', async ({
    page,
  }) => {
    const opfs = new OpfsPage(page);
    await seedDirectory(page, 'alpha/beta/gamma');
    await opfs.goto('alpha/beta/gamma/');

    await opfs.breadcrumb.getByText('beta', { exact: true }).click();

    expect(await opfs.breadcrumbSegments()).toEqual([
      '(Root)',
      'alpha',
      'beta',
    ]);
  });

  test('J3: full reload (F5) restores the path from the hash', async ({
    page,
  }) => {
    const opfs = new OpfsPage(page);
    await seedDirectory(page, 'persist/me');
    await opfs.goto('persist/me/');

    await page.reload();

    expect(await opfs.breadcrumbSegments()).toEqual([
      '(Root)',
      'persist',
      'me',
    ]);
  });
});
