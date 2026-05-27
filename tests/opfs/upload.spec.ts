// spec: specs/opfs-ui-test-plan.md
// TC-06: ローカルファイルのアップロード
// TC-07: アップロード — 同名ファイルの上書き確認
// TC-08: ファイルのアップロード前にボタンが無効

import { test, expect } from '@playwright/test';
import { gotoOPFS } from './helpers';

const UPLOAD_FILE = {
  name: 'hello.txt',
  mimeType: 'text/plain',
  buffer: Buffer.from('Hello, World!'),
};

test.describe('TC-06: ローカルファイルのアップロード', () => {
  test('ファイルが正常にアップロードされて一覧に表示される', async ({ page }) => {
    await gotoOPFS(page);

    // ファイルを選択する
    await page.setInputFiles('#upload-file', UPLOAD_FILE);

    // Upload as 欄にファイル名が自動入力される
    await expect(page.locator('#upload-name')).toHaveValue('hello.txt');

    // Upload ボタンが有効になる
    await expect(page.locator('#upload-upload')).toBeEnabled();

    // アップロード実行 — 成功アラートを受け取る
    const dialogPromise = page.waitForEvent('dialog');
    await page.click('#upload-upload');
    const dialog = await dialogPromise;
    expect(dialog.message()).toContain('"hello.txt"');
    await dialog.accept();

    // ファイル一覧に追加される
    const row = page.locator('.grid-row').filter({ hasText: 'hello.txt' });
    await expect(row).toBeVisible();
    await expect(row).toContainText('file');
    await expect(row).toContainText('13'); // "Hello, World!" = 13 bytes

    // 入力欄がクリアされる
    await expect(page.locator('#upload-name')).toHaveValue('');
  });
});

test.describe('TC-07: アップロード — 同名ファイルの上書き確認', () => {
  test('上書き確認ダイアログでキャンセルすると上書きされない', async ({ page }) => {
    await gotoOPFS(page);

    // 1回目のアップロード
    await page.setInputFiles('#upload-file', UPLOAD_FILE);
    const dialog1 = page.waitForEvent('dialog');
    await page.click('#upload-upload');
    await (await dialog1).accept();
    await expect(page.locator('.grid-row').filter({ hasText: 'hello.txt' })).toBeVisible();

    // 2回目: 同名ファイルをアップロード → 上書き確認ダイアログ（キャンセル）
    await page.setInputFiles('#upload-file', UPLOAD_FILE);
    await expect(page.locator('#upload-upload')).toBeEnabled();
    const dialog2 = page.waitForEvent('dialog');
    await page.click('#upload-upload');
    const dlg2 = await dialog2;
    expect(dlg2.message()).toContain('"hello.txt"');
    await dlg2.dismiss(); // キャンセル

    // ファイルは依然として存在する（行が 1 件のまま）
    await expect(page.locator('.grid-row')).toHaveCount(1);
  });

  test('上書き確認ダイアログで OK するとファイルが更新される', async ({ page }) => {
    await gotoOPFS(page);

    // 1回目のアップロード
    await page.setInputFiles('#upload-file', UPLOAD_FILE);
    const dialog1 = page.waitForEvent('dialog');
    await page.click('#upload-upload');
    await (await dialog1).accept();

    // 2回目: 同名ファイルをアップロード → 上書き確認ダイアログ（OK）
    const updated = { ...UPLOAD_FILE, buffer: Buffer.from('Updated content') };
    await page.setInputFiles('#upload-file', updated);
    const dialog2 = page.waitForEvent('dialog'); // 上書き confirm
    await page.click('#upload-upload');
    await (await dialog2).accept(); // OK

    const dialog3 = page.waitForEvent('dialog'); // 成功 alert
    const dlg3 = await dialog3;
    expect(dlg3.message()).toContain('"hello.txt"');
    await dlg3.accept();

    // 行が 1 件のまま（重複しない）
    await expect(page.locator('.grid-row')).toHaveCount(1);
  });
});

test.describe('TC-08: ファイル未選択時は Upload ボタンが無効', () => {
  test('ファイルを選択する前は Upload ボタンが disabled', async ({ page }) => {
    await gotoOPFS(page);
    await expect(page.locator('#upload-upload')).toBeDisabled();
  });
});
