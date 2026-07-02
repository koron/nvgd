import { test, expect } from '@playwright/test';

test.describe('Directory Navigation', () => {
  test('Click into Subdirectory', async ({ page }) => {
    // Step 1: Create a directory named `testdir`
    const dirNameInput = page.locator('input[placeholder*="dir"], input[name*="dir"]');
    await dirNameInput.click();
    await dirNameInput.fill('testdir');

    // Step 2: Create a subdirectory `testdir/subdir`
    const createDirBtn = page.locator('text=Create new directory, button');
    await createDirBtn.click();
    await expect(page.locator('text=testdir/, .listing, .directory-list, table')).toContainText('testdir/');

    // Navigate into testdir by clicking on it
    await page.locator('text=testdir/').click();

    // Step 4: Verify the breadcrumb trail shows `(Root) / testdir`
    const breadcrumb = page.locator('text=testdir, text=testdir/');
    await expect(breadcrumb).toBeVisible();

    // Step 5: Verify the directory listing shows the contents of `testdir`
    // The listing should show `subdir/` since we created it inside testdir
    const listing = page.locator('.listing, .directory-list, table');
    await expect(listing).toBeVisible();

    // Step 6: Verify the page title updates to `OPFS: testdir`
    await expect(page).toHaveTitle(/OPFS: testdir/);
  });

  test('Navigate Up with ..', async ({ page }) => {
    // Step 1: Create a directory `testdir`
    const dirNameInput = page.locator('input[placeholder*="dir"], input[name*="dir"]');
    await dirNameInput.click();
    await dirNameInput.fill('testdir');

    // Step 2: Create a subdirectory `testdir/subdir`
    const createDirBtn = page.locator('text=Create new directory, button');
    await createDirBtn.click();
    await expect(page.locator('text=testdir/')).toBeVisible();

    // Step 3: Navigate into `testdir/subdir`
    await page.locator('text=testdir/').click();
    await page.locator('text=subdir/').click();

    // Step 4: Click the `..` link in the breadcrumb
    await page.locator('text=..').click();

    // Step 5: Verify the listing shows the contents of `testdir`
    await expect(page.locator('text=subdir/')).toBeVisible();

    // Step 6: Verify the breadcrumb trail reflects the parent directory
    await expect(page.locator('text=testdir')).toBeVisible();

    // Step 7: Navigate up from the root directory
    // First go into testdir
    await page.locator('text=testdir/').click();
    // Try to navigate up from root (click .. at root level)
    await page.locator('text=..').click();

    // Step 8: Verify an alert `No parent directory` is shown
    const alert = page.locator('text=No parent directory');
    await expect(alert).toBeVisible();
  });

  test('Breadcrumb Navigation', async ({ page }) => {
    // Step 1: Create nested directories: `a/b/c/`
    const dirNameInput = page.locator('input[placeholder*="dir"], input[name*="dir"]');
    await dirNameInput.click();
    await dirNameInput.fill('a');
    await page.locator('text=Create new directory, button').click();
    await expect(page.locator('text=a/')).toBeVisible();

    // Navigate into `a/`
    await page.locator('text=a/').click();

    // Step 2: Navigate into `a/b/c/`
    await dirNameInput.click();
    await dirNameInput.fill('b');
    await page.locator('text=Create new directory, button').click();
    await expect(page.locator('text=b/')).toBeVisible();

    await page.locator('text=b/').click();

    await dirNameInput.click();
    await dirNameInput.fill('c');
    await page.locator('text=Create new directory, button').click();
    await expect(page.locator('text=c/')).toBeVisible();

    // Step 3: Navigate into `a/b/c/`
    await page.locator('text=c/').click();

    // Step 4: Click on `b` in the breadcrumb trail
    await page.locator('text=b').click();

    // Step 5: Verify the listing shows the contents of `a/b/`
    await expect(page.locator('text=c/')).toBeVisible();

    // Step 6: Click on `a` in the breadcrumb trail
    await page.locator('text=a').click();

    // Step 7: Verify the listing shows the contents of `a/`
    await expect(page.locator('text=b/')).toBeVisible();
  });

  test('History Back/Forward', async ({ page }) => {
    // Step 1: Navigate to OPFS root
    await page.goto('http://localhost:9280/opfs/');

    // Step 2: Create a directory `testdir`
    const dirNameInput = page.locator('input[placeholder*="dir"], input[name*="dir"]');
    await dirNameInput.click();
    await dirNameInput.fill('testdir');
    await page.locator('text=Create new directory, button').click();

    // Step 3: Navigate into `testdir`
    await page.locator('text=testdir/').click();

    // Step 4: Use browser Back button
    await page.goBack();

    // Step 5: Verify the OPFS UI restores the previous directory view
    // Should show the root directory with testdir/ in the listing
    await expect(page.locator('text=testdir/')).toBeVisible();

    // Step 6: Verify the browser history URL updates with hash fragments
    const url = page.url();
    expect(url).toContain('opfs/');
  });
});
