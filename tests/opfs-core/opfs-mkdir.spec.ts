import { test, expect } from '@playwright/test';

test.describe('Create Directory Operations', () => {
  test('Create Directory', async ({ page }) => {
    await page.goto('http://localhost:9280/opfs/');

    // Step 1: Enter `mydir` in the "Directory name to create" input
    const dirNameInput = page.locator('input[placeholder*="dir"], input[name*="dir"]');
    await dirNameInput.click();
    await dirNameInput.fill('mydir');

    // Step 2: Click "Create new directory" button
    await page.locator('text=Create new directory, button').click();

    // Step 3: Verify a new directory `mydir/` appears in the listing
    await expect(page.locator('text=mydir/')).toBeVisible();

    // Step 4: Click into `mydir` and verify it is empty
    await page.locator('text=mydir/').click();
    const listing = page.locator('.listing, .directory-list, table');
    await expect(listing).toBeVisible();

    // Step 5: Enter empty name and click "Create new directory"
    await dirNameInput.click();
    await dirNameInput.fill('');
    await page.locator('text=Create new directory, button').click();

    // Step 6: Verify an alert `Need directory name` is shown
    const alert = page.locator('text=Need directory name');
    await expect(alert).toBeVisible();

    // Step 7: Enter an existing directory name and click "Create new directory"
    await dirNameInput.click();
    await dirNameInput.fill('mydir');
    await page.locator('text=Create new directory, button').click();

    // Step 8: Verify an error is thrown
    const errorAlert = page.locator('text=Error, text=already exists, text=Error:');
    await expect(errorAlert).toBeVisible();
  });

  test('Create Nested Directories', async ({ page }) => {
    await page.goto('http://localhost:9280/opfs/');

    // Step 1: Create directory `project`
    const dirNameInput = page.locator('input[placeholder*="dir"], input[name*="dir"]');
    await dirNameInput.click();
    await dirNameInput.fill('project');
    await page.locator('text=Create new directory, button').click();
    await expect(page.locator('text=project/')).toBeVisible();

    // Step 2: Navigate into `project`
    await page.locator('text=project/').click();

    // Step 3: Create directory `src`
    await dirNameInput.click();
    await dirNameInput.fill('src');
    await page.locator('text=Create new directory, button').click();
    await expect(page.locator('text=src/')).toBeVisible();

    // Step 4: Navigate into `src`
    await page.locator('text=src/').click();

    // Step 5: Create directory `main`
    await dirNameInput.click();
    await dirNameInput.fill('main');
    await page.locator('text=Create new directory, button').click();
    await expect(page.locator('text=main/')).toBeVisible();

    // Step 6: Navigate back up and verify the full nested structure
    await page.locator('text=..').click();
    await expect(page.locator('text=src/')).toBeVisible();
    await page.locator('text=..').click();
    await expect(page.locator('text=project/')).toBeVisible();

    // Step 7: Verify breadcrumb trail correctly represents the nested path
    await expect(page.locator('text=project')).toBeVisible();
  });
});
