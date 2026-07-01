import { test, expect } from '@playwright/test';

test.describe('04-file-edit', () => {
  test('ファイル編集機能（64KiB未満）', async ({ page }) => {
    await page.goto('http://127.0.0.1:9280/opfs/');

    // Create a small file first
    await page.locator('#editor-name').fill('edit-test.txt');
    await page.locator('#editor-edit').fill('Original content for editing');
    await page.locator('#editor-save').click();
    await expect(page.locator('a:has-text("edit-test.txt")')).toBeVisible();

    // 1. 小ファイルを作成し、一覧で"Edit"リンクが表示されることを確認
    // expect: 64KiB未満のファイルに"Edit"アクションが表示される
    const editLink = page.locator('a:has-text("Edit")');
    await expect(editLink).toBeVisible();

    // 2. "Edit"リンクをクリックし、#editor-nameと#editor-editの値を確認
    await editLink.click();

    // expect: "Edit"クリックでエディタにファイル名と内容が読み込まれる
    await expect(page.locator('#editor-name')).toHaveValue('edit-test.txt');
    await expect(page.locator('#editor-edit')).toHaveValue('Original content for editing');

    // 3. #editor-editの値を変更し、#editor-saveをクリック
    await page.locator('#editor-edit').fill('Edited content!');
    await page.locator('#editor-save').click();

    // expect: エディタで内容を変更して保存できる

    // 4. ファイルを再度Editで読み込み、更新された内容を確認
    await editLink.click();

    // expect: 保存後、ファイルの内容が更新される
    await expect(page.locator('#editor-edit')).toHaveValue('Edited content!');
  });

  test('エディタのClear機能', async ({ page }) => {
    await page.goto('http://127.0.0.1:9280/opfs/');

    // 1. エディタに値を入力後、#editor-clearをクリック
    await page.locator('#editor-name').fill('clear-test.txt');
    await page.locator('#editor-edit').fill('Content to clear');
    await page.locator('#editor-clear').click();

    // expect: エディタに値が入力された状態で#editor-clearをクリックすると、両フィールドが空になる

    // 2. #editor-name.inputValue()が空であることを確認
    // expect: #editor-nameが空になる
    await expect(page.locator('#editor-name')).toHaveValue('');

    // 3. #editor-edit.inputValue()が空であることを確認
    // expect: #editor-editが空になる
    await expect(page.locator('#editor-edit')).toHaveValue('');
  });

  test('64KiB超のファイルにはEditが表示されない', async ({ page }) => {
    await page.goto('http://127.0.0.1:9280/opfs/');

    // Create a large file (>64KiB)
    const largeContent = 'x'.repeat(70 * 1024); // 70KiB

    // Upload large file
    await page.locator('#upload-file').setInputFiles({
      name: 'large-file.txt',
      contentType: 'text/plain',
      buffer: Buffer.from(largeContent),
    });

    // Handle alert for upload
    page.once('dialog', dialog => {
      dialog.accept();
    });

    await page.locator('#upload-upload').click();
    await expect(page.locator('a:has-text("large-file.txt")')).toBeVisible();

    // 1. 64KiB超のファイルを作成し、一覧で"Edit"リンクが表示されないことを確認
    // expect: 64KiB超のファイルには"Edit"リンクが表示されない
    const editLinks = page.locator('a:has-text("Edit")');
    const visibleEditLinks = await editLinks.all();
    let foundEdit = false;
    for (const link of visibleEditLinks) {
      const href = await link.getAttribute('href');
      if (href && href.includes('large-file.txt')) {
        foundEdit = true;
        break;
      }
    }
    expect(foundEdit).toBe(false);

    // 2. 一覧で"Save as"リンクが表示されることを確認
    // expect: 64KiB超のファイルには"Save as"リンクが表示される
    const saveAsLink = page.locator('a:has-text("Save as")');
    await expect(saveAsLink).toBeVisible();
  });
});
