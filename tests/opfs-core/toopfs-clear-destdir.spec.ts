import { test, expect } from '@playwright/test';

test.describe('toopfs Filter - Download Files from NVGD to OPFS', () => {
  test('Clear Destination Directory', async ({ page }) => {
    await page.goto('http://localhost:9280/file://?toopfs');

    // Enter a destination directory path
    const destDirInput = page.locator('input[placeholder*="OPFS directory"], input[name*="dest"]');
    await destDirInput.fill('test/output');

    // Click the "Clear" button next to the destination field
    const clearBtn = page.locator('text=Clear, button, :text("Clear")');
    await clearBtn.click();

    // Verify the destination directory field is cleared
    await expect(destDirInput).toHaveValue('');

    // Verify the field is focused after clearing
    await expect(destDirInput).toBeFocused();
  });

  test('Download Single File', async ({ page }) => {
    await page.goto('http://localhost:9280/file://?toopfs');

    // Select the file
    const checkbox = page.locator('input[type=checkbox], [class*="checkbox"]');
    await checkbox.first().check();

    // Set a destination directory
    const destDirInput = page.locator('input[placeholder*="OPFS directory"], input[name*="dest"]');
    await destDirInput.fill('single-output/');

    // Download the file
    const downloadBtn = page.locator('text=Download, button');
    await downloadBtn.click();
    await page.locator('text=Confirm, button, :text("Confirm")').click();

    // Verify the file appears in the OPFS destination directory
    await page.waitForEvent('page');
    const opfsPage = page.context().pages().find(p => p.url().includes('opfs'));
    expect(opfsPage).toBeTruthy();
  });
});
