import { test, expect } from '@playwright/test';

test.describe('File Upload Operations', () => {
  test('Upload File with Custom Name', async ({ page }) => {
    await page.goto('http://localhost:9280/opfs/');

    // Step 1: Select a local file `source.txt`
    const uploadInput = page.locator('input[type=file]');
    await uploadInput.setInputFiles('tests/testdata/sample.txt');

    // Step 2: Change the file name in the "File name to upload" field to `dest.txt`
    const nameField = page.locator('input[placeholder*="File name"]');
    await nameField.fill('dest.txt');

    // Step 3: Click "Upload" button
    await page.locator('text=Upload, button').click();

    // Step 4: Verify the file appears under its new name `dest.txt`
    await expect(page.locator('text=dest.txt')).toBeVisible();

    // Step 5: Verify the original file `source.txt` does not appear in the listing
    await expect(page.locator('text=source.txt')).not.toBeVisible();
  });
});
