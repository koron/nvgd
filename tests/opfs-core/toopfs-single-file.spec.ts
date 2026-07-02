import { test, expect } from '@playwright/test';

test.describe('toopfs Filter - Download Files from NVGD to OPFS', () => {
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
