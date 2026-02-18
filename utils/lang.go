package utils

import (
	"context"
	"fmt"
)

var Messages = map[string]map[string]string{
	"ar": {
		// ── اللغة ──
		"lang_confirmed": "لقد اخترت العربية يا %s!\nيمكنك الآن إرسال /help لعرض الأوامر.",

		// ── المساعدة ──
		"help_main":        "الأوامر المتاحة:",
		"btn_id_info":      "شرح أمر ID",
		"btn_carbon_info":  "شرح صنع الصور",
		"btn_replies_info": "إدارة الردود التلقائية",
		"back_btn":         "رجوع",

		// ── ID ──
		"id_description":  "أمر الـ ID يعرض معلوماتك.\n\n*طريقة الاستخدام:*\nأرسل `/id` لعرض معلوماتك.",
		"user_info_res":   "*معلومات المستخدم:*\n\nالاسم: %s\nالمعرف: %s\nالأيدي: `%d`\nاللغة: %s",
		"user_search_res": "تم البحث عن:\n`%s`",

		// ── Carbon ──
		"carbon_description": "أمر الكاربون يحول الكود إلى صورة احترافية.\n\n*طريقة الاستخدام:*\n`/carbon` متبوعاً بالكود.",
		"carbon_wait":        "جاري توليد الصورة...",
		"carbon_error":       "فشل توليد الصورة، تأكد من الكود.",

		// ── Info ──
		"info_msg": "معلومات البوت:\nاستخدم /help لعرض الأوامر.",

		// ── عام ──
		"not_admin":           "هذا الأمر للمشرفين فقط.",
		"operation_cancelled": "تم إلغاء العملية.",

		// ── addreply (خطوة بخطوة) ──
		"addreply_ask_trigger": "أرسل الكلمة التي تريد الرد عليها:\n\n_(أرسل *توقف* للإلغاء)_",
		"addreply_ask_reply":   "الكلمة: `%s`\n\nالآن أرسل الرد الذي سيظهر:\n\n_(أرسل *توقف* للإلغاء)_",
		"addreply_success":     "*تمت الإضافة!*\n\nعند إرسال: `%s`\nسيرد البوت بـ: %s",

		// ── delreply (خطوة بخطوة) ──
		"delreply_ask_trigger": "أرسل الكلمة التي تريد حذف ردها:\n\n_(أرسل *توقف* للإلغاء)_",
		"delreply_notfound":    "لا يوجد رد مسجل للكلمة: `%s`",
		"delreply_success":     "تم حذف الرد عن: `%s`",

		// ── listreplies ──
		"listreplies_empty":  "لا توجد ردود تلقائية حتى الآن.",
		"listreplies_header": "*الردود التلقائية المسجلة:*",

		// ── شرح الردود في Help ──
		"replies_description": "*الردود التلقائية:*\n\nتتيح لك إضافة ردود تلقائية على كلمات معينة.\n\n*الأوامر المتاحة (للمشرف فقط):*\n\n`/addreply` — إضافة رد جديد\nسيطلب منك الكلمة ثم الرد خطوة بخطوة.\n\n`/delreply` — حذف رد موجود\nسيطلب منك الكلمة التي تريد حذف ردها.\n\n`/listreplies` — عرض كل الردود المسجلة\n\n*مثال:*\nبعد إضافة رد على كلمة `هلو`،\nأي شخص يرسل `هلو` سيرد عليه البوت تلقائياً.",
	},

	"en": {
		// ── Language ──
		"lang_confirmed": "English chosen, %s!\nSend /help to see commands.",

		// ── Help ──
		"help_main":        "Available Commands:",
		"btn_id_info":      "ID Command Info",
		"btn_carbon_info":  "Carbon Image Info",
		"btn_replies_info": "Auto-Replies Management",
		"back_btn":         "Back",

		// ── ID ──
		"id_description":  "The ID command shows your info.\n\n*How to use:*\nSend `/id` to show your info.",
		"user_info_res":   "*User Info:*\n\nName: %s\nUsername: %s\nID: `%d`\nLanguage: %s",
		"user_search_res": "Search result for:\n`%s`",

		// ── Carbon ──
		"carbon_description": "Carbon converts your code into a professional image.\n\n*How to use:*\n`/carbon` followed by your code.",
		"carbon_wait":        "Generating image...",
		"carbon_error":       "Failed to generate image, check your code.",

		// ── Info ──
		"info_msg": "Bot Info:\nUse /help to see commands.",

		// ── General ──
		"not_admin":           "This command is for admins only.",
		"operation_cancelled": "Operation cancelled.",

		// ── addreply (step by step) ──
		"addreply_ask_trigger": "Send the word you want to auto-reply to:\n\n_(Send *stop* to cancel)_",
		"addreply_ask_reply":   "Word: `%s`\n\nNow send the reply text:\n\n_(Send *stop* to cancel)_",
		"addreply_success":     "*Added!*\n\nWhen someone sends: `%s`\nBot will reply: %s",

		// ── delreply (step by step) ──
		"delreply_ask_trigger": "Send the word you want to remove its reply:\n\n_(Send *stop* to cancel)_",
		"delreply_notfound":    "No reply found for: `%s`",
		"delreply_success":     "Reply deleted for: `%s`",

		// ── listreplies ──
		"listreplies_empty":  "No auto-replies registered yet.",
		"listreplies_header": "*Registered Auto-Replies:*",

		// ── Replies explanation in Help ──
		"replies_description": "*Auto-Replies:*\n\nAllows you to set automatic replies for specific words.\n\n*Commands (admins only):*\n\n`/addreply` — Add a new auto-reply\nThe bot will ask for the word then the reply step by step.\n\n`/delreply` — Delete an existing reply\nThe bot will ask which word to remove.\n\n`/listreplies` — Show all registered replies\n\n*Example:*\nAfter adding a reply for `hello`,\nanyone who sends `hello` will get an automatic reply.",
	},
}

func GetLang(userID int64) string {
	lang, err := RDB.Get(context.Background(), fmt.Sprintf("lang:%d", userID)).Result()
	if err != nil {
		return "ar"
	}
	return lang
}
