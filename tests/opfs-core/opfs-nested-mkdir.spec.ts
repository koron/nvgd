import { test, expect } from '@playwright/test';

test.describe('Create Directory Operations', () => {
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
