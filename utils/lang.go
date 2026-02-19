package utils

import (
	"context"
	"fmt"
)

var Messages = map[string]map[string]string{
	"ar": {
		"lang_confirmed": "لقد اخترت العربية يا %s!\nيمكنك الآن إرسال /help لعرض الأوامر.",

		"help_main":        "الأوامر المتاحة:",
		"btn_id_info":      "شرح أمر ID",
		"btn_carbon_info":  "شرح صنع الصور",
		"btn_replies_info": "إدارة الردود التلقائية",
		"btn_voices_info":  "استخدام الـ inline Use",
		"back_btn":         "رجوع",

		"id_description":  "أمر الـ ID يعرض معلوماتك.\n\n*طريقة الاستخدام:*\nأرسل `/id` لعرض معلوماتك.",
		"user_info_res":   "*معلومات المستخدم:*\n\nالاسم: %s\nالمعرف: %s\nالأيدي: `%d`\nاللغة: %s",
		"user_search_res": "تم البحث عن:\n`%s`",

		"carbon_description": "أمر الكاربون يحول الكود إلى صورة احترافية.\n\n*طريقة الاستخدام:*\n`/carbon` متبوعاً بالكود.",
		"carbon_wait":        " جاري توليد الصورة...",
		"carbon_error":       " فشل توليد الصورة، تأكد من الكود.",

		"info_msg": "معلومات البوت:\nاستخدم /help لعرض الأوامر.",

		"not_admin":           " هذا الأمر للمشرفين فقط.",
		"operation_cancelled": " تم إلغاء العملية.",

		"addreply_ask_trigger": " أرسل الكلمة التي تريد الرد عليها:\n\n_(أرسل *توقف* للإلغاء)_",
		"addreply_ask_reply":   " الكلمة: `%s`\n\nالآن أرسل الرد الذي سيظهر:\n\n_(أرسل *توقف* للإلغاء)_",
		"addreply_success":     " *تمت الإضافة!*\n\nعند إرسال: `%s`\nسيرد البوت بـ: %s",

		"delreply_ask_trigger": " أرسل الكلمة التي تريد حذف ردها:\n\n_(أرسل *توقف* للإلغاء)_",
		"delreply_notfound":    " لا يوجد رد مسجل للكلمة: `%s`",
		"delreply_success":     " تم حذف الرد عن: `%s`",

		"listreplies_empty":  " لا توجد ردود تلقائية حتى الآن.",
		"listreplies_header": " *الردود التلقائية المسجلة:*",

		"replies_description": " *الردود التلقائية:*\n\nتتيح لك إضافة ردود تلقائية على كلمات معينة.\n\n*الأوامر المتاحة (للمشرف فقط):*\n\n`/addreply` — إضافة رد جديد\nسيطلب منك الكلمة ثم الرد خطوة بخطوة.\n\n`/delreply` — حذف رد موجود\nسيطلب منك الكلمة التي تريد حذف ردها.\n\n`/listreplies` — عرض كل الردود المسجلة\n\n*مثال:*\nبعد إضافة رد على كلمة `هلو`،\nأي شخص يرسل `هلو` سيرد عليه البوت تلقائياً.",

		"addvoice_reply_required": " يجب الرد على رسالة صوتية!\n\nارسل `/addvoice الاسم` كرد على صوت.",
		"addvoice_voice_required": " الرسالة التي رددت عليها ليست صوتاً!",
		"addvoice_name_required":  " يجب كتابة الاسم!\n\nمثال: `/addvoice احمد`",
		"addvoice_success":        " تمت إضافة صوت *%s* بنجاح!\n\nيمكن البحث عنه الآن عبر الـ inline.",
		"delvoice_usage":          " يجب كتابة الاسم!\n\nمثال: `/delvoice احمد`",
		"delvoice_notfound":       " لا يوجد صوت مسجل باسم: `%s`",
		"delvoice_success":        " تم حذف صوت: `%s`",
		"listvoices_empty":        " لا توجد أصوات مسجلة حتى الآن.",
		"listvoices_header":       " *الأصوات المسجلة:*",
		"voices_description":      " *البحث بالـ Inline:*\n\nتتيح لك هذه الميزة البحث وإرسال البصمات الصوتية وملفات الـ IPA مباشرة.\n\n*البصمات الصوتية:*\nاكتب `@معرف_البوت aud` متبوعاً باسم البصمة.\n\n*ملفات الـ IPA:*\nاكتب `@معرف_البوت ipa` متبوعاً باسم الملف.\n\n*الأوامر (للمشرف فقط):*\n`/addvoice الاسم` — حفظ بصمة\n`/addipa الاسم` — حفظ ملف IPA",

		"addipa_ask_doc":        "الكلمة: `%s`\n\nالآن أرسل ملف الـ IPA:\n\n_(أرسل *توقف* للإلغاء)_",
		"addipa_reply_required": "يجب الرد على ملف IPA أو إرساله مع التعليق!",
		"addipa_doc_required":   "الرسالة ليست ملفا!",
		"addipa_ipa_required":   "الملف ليس IPA! يجب ان ينتهي بـ .ipa",
		"addipa_name_required":  "يجب كتابة الاسم!\n\nمثال: `/addipa MyApp`",
		"addipa_success":        "تمت اضافة *%s* بنجاح!\n\nيمكن البحث عنه الآن عبر الـ inline بدون الحاجة لكتابة ipa.",
		"delipa_usage":          "يجب كتابة الاسم!\n\nمثال: `/delipa MyApp`",
		"delipa_notfound":       "لا يوجد ملف مسجل باسم: `%s`",
		"delipa_success":        "تم حذف: `%s`",
		"listipa_empty":         "لا توجد ملفات IPA مسجلة حتى الان.",
		"listipa_header":        "*ملفات IPA المسجلة:*",
		"inline_help_title":     "كيفية استخدام البحث؟",
		"inline_help_desc":      "يجب كتابة كلمة aud أو ipa للبحث.",
		"inline_help_text":      "🔍 *طريقة استخدام البحث:*\n\nيجب كتابة الكلمة المفتاحية قبل البحث:\n\n1. للبحث عن بصمة: `@bot aud الاسم`\n2. للبحث عن ملف: `@bot ipa الاسم`\n\n*ملاحظة:* لن تظهر أي نتائج إذا لم تكتب الكلمة المفتاحية أولاً.",
	},

	"en": {
		"lang_confirmed": "English chosen, %s!\nSend /help to see commands.",

		"help_main":        "Available Commands:",
		"btn_id_info":      "ID Command Info",
		"btn_carbon_info":  "Carbon Image Info",
		"btn_replies_info": "Auto-Replies Management",
		"btn_voices_info":  "Inline Use",
		"back_btn":         "Back",

		"id_description":  "The ID command shows your info.\n\n*How to use:*\nSend `/id` to show your info.",
		"user_info_res":   "*User Info:*\n\nName: %s\nUsername: %s\nID: `%d`\nLanguage: %s",
		"user_search_res": "Search result for:\n`%s`",

		"carbon_description": "Carbon converts your code into a professional image.\n\n*How to use:*\n`/carbon` followed by your code.",
		"carbon_wait":        " Generating image...",
		"carbon_error":       " Failed to generate image, check your code.",

		"info_msg": "Bot Info:\nUse /help to see commands.",

		"not_admin":           " This command is for admins only.",
		"operation_cancelled": " Operation cancelled.",

		"addreply_ask_trigger": " Send the word you want to auto-reply to:\n\n_(Send *stop* to cancel)_",
		"addreply_ask_reply":   " Word: `%s`\n\nNow send the reply text:\n\n_(Send *stop* to cancel)_",
		"addreply_success":     " *Added!*\n\nWhen someone sends: `%s`\nBot will reply: %s",

		"delreply_ask_trigger": " Send the word you want to remove its reply:\n\n_(Send *stop* to cancel)_",
		"delreply_notfound":    " No reply found for: `%s`",
		"delreply_success":     " Reply deleted for: `%s`",

		"listreplies_empty":  " No auto-replies registered yet.",
		"listreplies_header": " *Registered Auto-Replies:*",

		"replies_description": " *Auto-Replies:*\n\nAllows you to set automatic replies for specific words.\n\n*Commands (admins only):*\n\n`/addreply` — Add a new auto-reply\nThe bot will ask for the word then the reply step by step.\n\n`/delreply` — Delete an existing reply\nThe bot will ask which word to remove.\n\n`/listreplies` — Show all registered replies\n\n*Example:*\nAfter adding a reply for `hello`,\nanyone who sends `hello` will get an automatic reply.",

		"addvoice_reply_required": " You must reply to a voice message!\n\nSend `/addvoice name` as a reply to a voice.",
		"addvoice_voice_required": " The message you replied to is not a voice!",
		"addvoice_name_required":  " You must provide a name!\n\nExample: `/addvoice ahmed`",
		"addvoice_success":        " Voice *%s* added successfully!\n\nYou can now search for it via inline.",
		"delvoice_usage":          " You must provide a name!\n\nExample: `/delvoice ahmed`",
		"delvoice_notfound":       " No voice found with name: `%s`",
		"delvoice_success":        " Voice deleted: `%s`",
		"listvoices_empty":        " No voices registered yet.",
		"listvoices_header":       " *Registered Voices:*",
		"voices_description":      " *Inline Search:*\n\nSearch and send voice fingerprints and IPA files directly.\n\n*Voices:*\nType `@bot aud` followed by the name.\n\n*IPAs:*\nType `@bot ipa` followed by the name.\n\n*Commands (Admins only):*\n`/addvoice Name` — Save a voice\n`/addipa Name` — Save an IPA",

		"addipa_ask_doc":        "Name: `%s`\n\nNow send the IPA file:\n\n_(Send *stop* to cancel)_",
		"addipa_reply_required": "You must reply to an IPA file or send it as a caption!",
		"addipa_doc_required":   "The message is not a file!",
		"addipa_ipa_required":   "The file is not an IPA! Must end with .ipa",
		"addipa_name_required":  "You must provide a name!\n\nExample: `/addipa MyApp`",
		"addipa_success":        "*%s* added successfully!\n\nYou can search for it via inline without needing 'ipa' prefix.",
		"delipa_usage":          "You must provide a name!\n\nExample: `/delipa MyApp`",
		"delipa_notfound":       "No file found with name: `%s`",
		"delipa_success":        "Deleted: `%s`",
		"listipa_empty":         "No IPA files registered yet.",
		"listipa_header":        "*Registered IPA Files:*",
		"inline_help_title":     "How to use search?",
		"inline_help_desc":      "You must type 'aud' or 'ipa' to search.",
		"inline_help_text":      "🔍 *How to use search:*\n\nKeywords are required before the search term:\n\n1. For voices: `@bot aud Name`\n2. For IPAs: `@bot ipa Name`\n\n*Note:* No results will appear if you don't type the keyword first.",
	},
}

func GetLang(userID int64) string {
	lang, err := RDB.Get(context.Background(), fmt.Sprintf("lang:%d", userID)).Result()
	if err != nil {
		return "ar"
	}
	return lang
}
