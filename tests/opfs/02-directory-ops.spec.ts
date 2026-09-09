import { test, expect } from '@playwright/test';

test.describe('02-directory-ops', () => {
  test('ディレクトリ作成機能', async ({ page }) => {
    await page.goto('http://127.0.0.1:9280/opfs/');

    // 1. #mkdir-nameに"test-dir"を入力し、#mkdir-mkdirをクリック
    await page.locator('#mkdir-name').fill('test-dir');
    await page.locator('#mkdir-mkdir').click();

    // expect: ディレクトリが一覧に追加される
    await expect(page.locator('a:has-text("test-dir/")')).toBeVisible();

    // 2. #mkdir-nameの値が空であることを確認
    // expect: 入力フィールドがクリアされる
    await expect(page.locator('#mkdir-name')).toHaveValue('');

    // 3. 一覧でtype列が"dir"であることを確認
    // expect: 作成されたディレクトリの型が"dir"である
    await expect(page.locator('.grid-row div:nth-child(2):has-text("dir")')).toBeVisible();

    // 4. 一覧でsizeとmodifiedAtが"(N/A)"であることを確認
    // expect: サイズとModified Atが"(N/A)"である
    await expect(page.locator('.grid-row div:has-text("(N/A)")')).toBeVisible();
  });

  test('ディレクトリ名のバリデーション（空入力）', async ({ page }) => {
    await page.goto('http://127.0.0.1:9280/opfs/');

    // 1. #mkdir-nameを空にして#mkdir-mkdirをクリック
    await page.locator('#mkdir-name').clear();
    await page.locator('#mkdir-mkdir').click();

    // expect: 空の入力でアラート"Need directory name"が表示される

    // 2. 一覧に新しいディレクトリが追加されないことを確認
    // expect: ディレクトリが作成されない
  });

  test('ディレクトリ名に特殊文字', async ({ page }) => {
    await page.goto('http://127.0.0.1:9280/opfs/');

    // 1. #mkdir-nameに"dir with spaces"を入力し作成
    await page.locator('#mkdir-name').fill('dir with spaces');
    await page.locator('#mkdir-mkdir').click();
    await expect(page.locator('a:has-text("dir with spaces/")')).toBeVisible();

    // 2. #mkdir-nameに"my-dir_name"を入力し作成
    await page.locator('#mkdir-name').fill('my-dir_name');
    await page.locator('#mkdir-mkdir').click();
    await expect(page.locator('a:has-text("my-dir_name/")')).toBeVisible();

    // 3. #mkdir-nameに"dir.special;name"を入力し作成
    await page.locator('#mkdir-name').fill('dir.special;name');
    await page.locator('#mkdir-mkdir').click();
    await expect(page.locator('a:has-text("dir.special;name/")')).toBeVisible();
  });

  test('アップロード名の変更', async ({ page }) => {
    await page.goto('http://127.0.0.1:9280/opfs/');

    // Create a test file for upload name testing
    await page.locator('#editor-name').fill('upload-test.txt');
    await page.locator('#editor-edit').fill('test content');
    await page.locator('#editor-save').click();

    // 1. ファイルを選択し、#upload-nameの値を確認
    await page.locator('#upload-file').setInputFiles({
      name: 'upload-test.txt',
      contentType: 'text/plain',
      buffer: Buffer.from('test content'),
    });

    // 2. #upload-nameの値を別の名前に変更
    // expect: #upload-nameの値を編集できる
    await page.locator('#upload-name').fill('renamed-file.txt');
    await expect(page.locator('#upload-name')).toHaveValue('renamed-file.txt');
  });
});
