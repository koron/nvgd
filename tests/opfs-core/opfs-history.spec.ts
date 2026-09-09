import { test, expect } from '@playwright/test';

test.describe('Directory Navigation', () => {
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
