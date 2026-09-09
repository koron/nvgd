import { test, expect } from '@playwright/test';

test.describe('06-file-delete', () => {
  test('単体ファイル削除', async ({ page }) => {
    await page.goto('http://127.0.0.1:9280/opfs/');

    // Create a file first
    await page.locator('#editor-name').fill('delete-single.txt');
    await page.locator('#editor-edit').fill('Content for single delete');
    await page.locator('#editor-save').click();
    await expect(page.locator('a:has-text("delete-single.txt")')).toBeVisible();

    // 1. ファイルのチェックボックスをチェックし、#command-deleteのdisabled状態を確認
    await page.locator('input[name="delete-single.txt"]').check();

    // expect: ファイルのチェックボックスを選択後、Deleteボタンが有効になる
    await expect(page.locator('#command-delete')).not.toBeDisabled();

    // 2. #command-deleteをクリックし、confirmダイアログが表示されることを確認
    // Handle confirm dialog - first dismiss to test dismiss behavior
    let firstDialog = true;
    page.once('dialog', async dialog => {
      expect(dialog.message()).toContain('delete-single.txt');
      if (firstDialog) {
        firstDialog = false;
        await dialog.dismiss();
      }
    });

    await page.locator('#command-delete').click();

    // expect: Deleteクリックで確認ダイアログが表示される

    // 3. confirmダイアログでacceptし、一覧からファイルが消えることを確認
    // The file should still be there since we dismissed
    await expect(page.locator('a:has-text("delete-single.txt")')).toBeVisible();

    // Re-check the checkbox
    await page.locator('input[name="delete-single.txt"]').check();

    // Handle confirm dialog - accept this time
    page.once('dialog', async dialog => {
      expect(dialog.message()).toContain('delete-single.txt');
      await dialog.accept();
    });

    await page.locator('#command-delete').click();

    // expect: 確認でOKを選択するとファイルが削除される
    await expect(page.locator('a:has-text("delete-single.txt")')).not.toBeVisible();

    // 4. confirmダイアログでdismissし、一覧に残ることを確認
    // Already tested above with dismiss
  });

  test('複数ファイルの削除', async ({ page }) => {
    await page.goto('http://127.0.0.1:9280/opfs/');

    // Create multiple files
    await page.locator('#editor-name').fill('multi-delete-1.txt');
    await page.locator('#editor-edit').fill('Content 1');
    await page.locator('#editor-save').click();

    await page.locator('#editor-name').fill('multi-delete-2.txt');
    await page.locator('#editor-edit').fill('Content 2');
    await page.locator('#editor-save').click();

    await expect(page.locator('a:has-text("multi-delete-1.txt")')).toBeVisible();
    await expect(page.locator('a:has-text("multi-delete-2.txt")')).toBeVisible();

    // 1. 複数のファイルをチェック
    await page.locator('input[name="multi-delete-1.txt"]').check();
    await page.locator('input[name="multi-delete-2.txt"]').check();

    // expect: 複数のファイルを選択後、Deleteボタンが有効になる
    await expect(page.locator('#command-delete')).not.toBeDisabled();

    // 2. confirmダイアログのメッセージに全ファイル名が含まれることを確認
    // Handle confirm dialog
    page.once('dialog', async dialog => {
      const message = dialog.message();
      expect(message).toContain('multi-delete-1.txt');
      expect(message).toContain('multi-delete-2.txt');
      await dialog.accept();
    });

    await page.locator('#command-delete').click();

    // expect: 確認ダイアログに選択したファイル名が全て表示される

    // 3. confirmでacceptし、全ファイルが消えることを確認
    await expect(page.locator('a:has-text("multi-delete-1.txt")')).not.toBeVisible();
    await expect(page.locator('a:has-text("multi-delete-2.txt")')).not.toBeVisible();

    // expect: 確認でOKを選択すると全ファイルが削除される
  });

  test('ディレクトリ選択時の再帰削除', async ({ page }) => {
    await page.goto('http://127.0.0.1:9280/opfs/');

    // Create a directory with a file inside
    await page.locator('#mkdir-name').fill('delete-dir');
    await page.locator('#mkdir-mkdir').click();

    // Navigate into the directory
    await page.locator('a:has-text("delete-dir/")').click();
    await expect(page).toHaveTitle('OPFS: /delete-dir/');

    // Create a file inside
    await page.locator('#editor-name').fill('delete-dir-file.txt');
    await page.locator('#editor-edit').fill('File inside directory');
    await page.locator('#editor-save').click();
    await expect(page.locator('a:has-text("delete-dir-file.txt")')).toBeVisible();

    // Navigate back to root
    await page.locator('#header span').click();
    await expect(page).toHaveTitle('OPFS: /');

    // 1. サブディレクトリ（中にファイルあり）を選択し削除
    await page.locator('input[name="delete-dir/"]').check();
    await expect(page.locator('#command-delete')).not.toBeDisabled();

    // Handle confirm dialog
    page.once('dialog', async dialog => {
      expect(dialog.message()).toContain('delete-dir');
      await dialog.accept();
    });

    await page.locator('#command-delete').click();

    // expect: ディレクトリを選択してDeleteすると、内容を含むディレクトリが削除される
    await expect(page.locator('a:has-text("delete-dir/")')).not.toBeVisible();
    await expect(page.locator('a:has-text("delete-dir-file.txt")')).not.toBeVisible();

    // 2. ディレクトリのcheckbox nameが"dir/"形式であることを確認
    // Already verified above - the checkbox name was "delete-dir/"
  });
});
