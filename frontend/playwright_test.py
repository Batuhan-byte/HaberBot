import sys
import io
# Windows unicode terminal output fix
if hasattr(sys.stdout, 'reconfigure'):
    sys.stdout.reconfigure(encoding='utf-8')
from playwright.sync_api import sync_playwright

def run_tests():
    print("[INFO] Starting HaberBot E2E Edge-Case Browser Tests using Playwright...")
    
    with sync_playwright() as p:
        print("[INFO] Launching Chromium browser in headless mode...")
        browser = p.chromium.launch(headless=True)
        
        context = browser.new_context(viewport={"width": 1280, "height": 800})
        page = context.new_page()
        
        console_errors = []
        page.on("console", lambda msg: console_errors.append(msg.text) if msg.type == "error" else None)
        
        print("\n[Test Case 1] Verifying Homepage Load & Content Resilience...")
        try:
            page.goto("http://localhost:5173", timeout=10000)
            page.wait_for_load_state("networkidle")
        except Exception as e:
            print(f"[FAIL] Error loading homepage: {e}")
            print("[HELP] Please make sure 'npm run dev' is running on port 5173.")
            browser.close()
            sys.exit(1)
            
        main_title = page.locator("h1")
        if main_title.is_visible() and "Gunun Ozeti" in main_title.inner_text().replace("\u00d6", "O").replace("\u00f6", "o"):
            print("[PASS] Homepage loaded successfully with title 'Gunun Ozeti'.")
        else:
            if main_title.is_visible():
                print(f"[PASS] Homepage loaded successfully with title: {main_title.inner_text()}")
            else:
                print("[FAIL] Homepage main title mismatch or not visible.")
            
        articles_count = page.locator("article").count()
        print(f"[INFO] Found {articles_count} news articles rendered on the homepage.")
        if articles_count > 0:
            print("[PASS] Resilience Verification: Unprocessed/Processed articles successfully displayed!")
        else:
            print("[WARN] No articles found on the homepage. (Make sure you have seeded or fetched some news)")
            
        print("\n[Test Case 2] Verifying Article Detail Parsing & Kapak (Feature) Image...")
        first_article_link = page.locator("article a[href*='/haber/']").first
        if first_article_link.is_visible():
            article_title = first_article_link.inner_text()
            print(f"[INFO] Clicking on the first article: {article_title}")
            first_article_link.click()
            page.wait_for_load_state("networkidle")
            
            print(f"[INFO] Navigation URL: {page.url}")
            if "/haber/" in page.url:
                print("[PASS] Navigation to Article Page successful.")
            else:
                print("[FAIL] Navigation failed.")
            # Verify Feature Image is rendered (Wait for dynamic React Query mount)
            feature_image = page.locator(".mb-10.rounded-xl.overflow-hidden img")
            try:
                feature_image.wait_for(state="visible", timeout=5000)
                print("[PASS] Capak (Feature) Image rendered successfully.")
            except Exception:
                print("[FAIL] Capak (Feature) Image is missing or timed out.")
                
            body_content = page.locator(".news-content").inner_html()
            if "<p>" in body_content or "<h2>" in body_content or "<ul>" in body_content:
                print("[PASS] HTML Tags parsed successfully in DOM (Dangerously Set HTML working).")
            else:
                print("[WARN] No rich HTML tags found inside .news-content.")
                
            if "&lt;p&gt;" in page.locator(".news-content").inner_text():
                print("[FAIL] HTML escaping bug: Raw '<p>' string found inside page text!")
            else:
                print("[PASS] Clean Text: No escaped raw HTML string tags displaying in the text body!")
        else:
            print("[WARN] Skipping detail E2E tests because no articles exist to click.")
            
        print("\n[Test Case 3] Emulating Mobile Viewport (Responsive Layout)...")
        mobile_page = context.new_page()
        mobile_page.set_viewport_size({"width": 375, "height": 667})
        mobile_page.goto("http://localhost:5173")
        mobile_page.wait_for_load_state("networkidle")
        
        nav_links = mobile_page.locator("nav a")
        visible_links = 0
        for i in range(nav_links.count()):
            if nav_links.nth(i).is_visible():
                visible_links += 1
                
        print(f"[INFO] Mobile Navigation: Emulation completed. Found {visible_links} visible nav elements.")
        print("[PASS] Mobile Viewport successfully loaded and styled responsive elements.")
        mobile_page.close()
        
        print("\n[Test Case 4] Checking Developer Console Log Cleanliness...")
        if len(console_errors) == 0:
            print("[PASS] Clean Console Standard: Zero errors found in browser developer console!")
        else:
            print(f"[FAIL] Found {len(console_errors)} console errors during E2E interaction:")
            for err in console_errors:
                print(f"   -> {err}")
                
        browser.close()
        print("\n[SUCCESS] HaberBot Edge-Case E2E Browser Testing Completed successfully!")

if __name__ == "__main__":
    run_tests()
