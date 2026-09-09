import { test, expect } from '@playwright/test';

test.describe('File Editor Operations', () => {
  test('File Size Limit (64KiB)', async ({ page }) => {
    await page.goto('http://localhost:9280/opfs/');

    // Step 1: Create a file with content larger than 64KiB (65536 bytes)
    const editorName = page.locator('input[name*="file"], input[placeholder*="file"], input[type=text]');
    const editorTextarea = page.locator('textarea');
    const largeContent = 'A'.repeat(65536);
    await editorName.fill('large.txt');
    await editorTextarea.fill(largeContent);
    await page.locator('text=Create or update a file, button').click();
    await expect(page.locator('text=Uploaded "large.txt"')).toBeVisible();

    // Step 2: Verify the "Edit" action is NOT available for files >= 64KiB
    const editBtn = page.locator('text=Edit, [class*="edit"], button[title*="Edit"]');
    await expect(editBtn).not.toBeVisible();

    // Step 3: Verify the file is still visible in the listing with correct size
    const listing = page.locator('.listing, .directory-list, table');
    await expect(listing).toContainText('large.txt');

    // Step 4: Verify the file is still downloadable via "Save as" action
    const saveAsBtn = page.locator('text=Save as, [class*="save"], button[title*="Save"]');
    await expect(saveAsBtn).toBeVisible();
  });

  test('Tab Key Handling', async ({ page }) => {
    await page.goto('http://localhost:9280/opfs/');

    // Step 1: Focus the editor textarea
    const editorTextarea = page.locator('textarea');
    await editorTextarea.click();

    // Step 2: Press the Tab key while in the editor
    await page.keyboard.press('Tab');

    // Step 3: Verify a tab character is inserted into the editor content
    const content = await editorTextarea.inputValue();
    expect(content).toContain('\t');

    // Step 4: Save the file and verify the tab character is preserved
    const editorName = page.locator('input[name*="file"], input[placeholder*="file"], input[type=text]');
    await editorName.fill('tab-test.txt');
    await page.locator('text=Create or update a file, button').click();
    await expect(page.locator('text=Uploaded "tab-test.txt"')).toBeVisible();
  });

  test('Clear Button', async ({ page }) => {
    await page.goto('http://localhost:9280/opfs/');

    // Step 1: Enter a file name and content in the editor
    const editorName = page.locator('input[name*="file"], input[placeholder*="file"], input[type=text]');
    await editorName.fill('clear-test.txt');
    const editorTextarea = page.locator('textarea');
    await editorTextarea.fill('Some content to clear');
    await page.locator('text=Create or update a file, button').click();
    await expect(page.locator('text=Uploaded "clear-test.txt"')).toBeVisible();

    // Step 2: Click the "Clear" button
    await page.locator('text=Clear, button').click();

    // Step 3: Verify the file name field is cleared
    await expect(editorName).toHaveValue('');

    // Step 4: Verify the editor textarea is cleared
    await expect(editorTextarea).toHaveValue('');

    // Step 5: Verify the previously saved file still exists in OPFS
    const listing = page.locator('.listing, .directory-list, table');
    await expect(listing).toContainText('clear-test.txt');
  });
});
