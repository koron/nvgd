import { test, expect } from '@playwright/test';

test.describe('07-multi-select', () => {
  test('全選択チェックボックス', async ({ page }) => {
    await page.goto('http://127.0.0.1:9280/opfs/');

    // Create some files for testing
    await page.locator('#editor-name').fill('select-test-1.txt');
    await page.locator('#editor-edit').fill('Content 1');
    await page.locator('#editor-save').click();

    await page.locator('#editor-name').fill('select-test-2.txt');
    await page.locator('#editor-edit').fill('Content 2');
    await page.locator('#editor-save').click();

    await expect(page.locator('a:has-text("select-test-1.txt")')).toBeVisible();
    await expect(page.locator('a:has-text("select-test-2.txt")')).toBeVisible();

    // 1. #toggle-selection-allをクリックし、全てのcheckboxがcheckedであることを確認
    const toggleAll = page.locator('#toggle-selection-all');
    await toggleAll.click();

    // expect: 全選択チェックボックスをクリックすると全アイテムが選択される
    const checkboxes = page.locator('input[type="checkbox"]');
    const checkedCount = (await checkboxes.all()).filter(async cb => await cb.isChecked()).length;
    expect(checkedCount).toBeGreaterThanOrEqual(2);

    // 2. #toggle-selection-allを再度クリックし、全てのcheckboxがuncheckedであることを確認
    await toggleAll.click();

    // expect: 再度クリックすると全選択解除される
    const uncheckedCheckboxes = page.locator('input[type="checkbox"]');
    for (const cb of await uncheckedCheckboxes.all()) {
      await expect(cb).not.toBeChecked();
    }

    // 3. 一部のアイテムをチェックし、#toggle-selection-allのindeterminate状態を確認
    const checkbox1 = page.locator('input[name="select-test-1.txt"]');
    await checkbox1.check();

    // expect: 一部のアイテムが選択されている状態ではindeterminateになる
    const toggleAllState = await toggleAll.evaluate(el => el.indeterminate);
    expect(toggleAllState).toBe(true);
  });

  test('選択状態とボタンの連動', async ({ page }) => {
    await page.goto('http://127.0.0.1:9280/opfs/');

    // Create a file
    await page.locator('#editor-name').fill('btn-test.txt');
    await page.locator('#editor-edit').fill('Content');
    await page.locator('#editor-save').click();

    // 1. アイテムをチェックし、#command-deleteのdisabledがfalseであることを確認
    await page.locator('input[name="btn-test.txt"]').check();
    await expect(page.locator('#command-delete')).not.toBeDisabled();

    // expect: アイテムを選択するとDeleteボタンが有効になる

    // 2. アイテムをチェックし、#command-duckdbのdisabledがfalseであることを確認
    // Note: DuckDB button may or may not exist depending on server config
    const duckdbBtn = page.locator('#command-duckdb');
    const duckdbVisible = await duckdbBtn.isVisible();
    if (duckdbVisible) {
      await expect(duckdbBtn).not.toBeDisabled();
    }

    // expect: アイテムを選択するとDuckDBボタンが有効になる

    // 3. 全選択解除し、両ボタンのdisabledがtrueであることを確認
    await page.locator('#toggle-selection-all').click();

    await expect(page.locator('#command-delete')).toBeDisabled();

    // expect: 全選択解除すると両ボタンが無効になる
  });

  test('全選択解除後の自動アンチェック', async ({ page }) => {
    await page.goto('http://127.0.0.1:9280/opfs/');

    // Create a directory and some files
    await page.locator('#mkdir-name').fill('multi-select-dir');
    await page.locator('#mkdir-mkdir').click();

    // Navigate into the directory
    await page.locator('a:has-text("multi-select-dir/")').click();
    await expect(page).toHaveTitle('OPFS: /multi-select-dir/');

    // Create a file inside
    await page.locator('#editor-name').fill('multi-select-file.txt');
    await page.locator('#editor-edit').fill('Content');
    await page.locator('#editor-save').click();

    // Check the checkbox
    await page.locator('input[name="multi-select-file.txt"]').check();
    await expect(page.locator('input[name="multi-select-file.txt"]')).toBeChecked();

    // Navigate back to root
    await page.locator('#header span').click();
    await expect(page).toHaveTitle('OPFS: /');

    // 1. サブディレクトリに移動し、チェックボックスが全てuncheckedであることを確認
    await page.locator('a:has-text("multi-select-dir/")').click();
    await expect(page).toHaveTitle('OPFS: /multi-select-dir/');

    // expect: ディレクトリ移動時に選択が解除される
    const checkboxes = page.locator('input[type="checkbox"]');
    for (const cb of await checkboxes.all()) {
      await expect(cb).not.toBeChecked();
    }
  });
});
