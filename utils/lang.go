package utils

import (
	"context"
	"fmt"
)

var Messages = map[string]map[string]string{
	"ar": {
		"lang_confirmed": "لقد اخترت العربية يا %s! \nيمكنك الآن إرسال /help لعرض الأوامر.",
		"help_main":      "الأوامر المتاحة:",
		"btn_id_info":    "شرح أمر ID",
		"id_description": "أمر الـ ID يعرض معلوماتك الخاصة، ويمكنك البحث عن مستخدمين آخرين عبر الـ ID أو اسم المستخدم.\n\n**طريقة الاستخدام:**\nأرسل `/id` لعرض معلوماتك، أو أرسل `/id 123456` للبحث عن شخص آخر.",
		"user_info_res": "**معلومات المستخدم:**\n\nالاسم: %s\nالمعرف: %s\nالأيدي: `%d`\nاللغة: %s",
		"user_search_res": "تم البحث عن:\n`%s`\n\nلا يمكن جلب معلومات المستخدم إلا إذا كان في نفس الكروب أو تفاعل مع البوت.",
		"btn_carbon_info": "شرح صنع الصور ",
		"carbon_description": "أمر الكاربون يحول الكود البرمجي إلى صورة احترافية.\n\n**طريقة الاستخدام:**\nأرسل `/carbon` متبوعاً بالكود الخاص بك.\n\nمثال:\n`/carbon fmt.Println(\"Hello\")`",
		"carbon_wait": "جاري توليد الصورة... قد يستغرق الأمر ثوانٍ.",
		"carbon_error": "فشل توليد الصورة، تأكد من صحة الكود.",
		"back_btn": "رجوع",

	},

	"en": {
		"lang_confirmed": "English chosen, %s!\nSend /help to see commands.",
		"help_main":      "Available Commands:",
		"btn_id_info":    "ID Command Info",
		"id_description": "The ID command shows your information and allows you to search for a user by ID or username.\n\n*How to use:*\nSend `/id` to show your info, or send `/id 123456` to search for someone else.",
		"user_info_res": "**User Info:**\n\nName: %s\nUsername: %s\nID: `%d`\nLanguage: %s",
		"user_search_res": "Search result for:\n`%s`\n\nUser info can only be retrieved if they are in the same chat or interacted with the bot.",
		"btn_carbon_info": "Carbon Image Info",
		"carbon_description": "Carbon command converts your code into a professional image.\n\n**How to use:**\nSend `/carbon` followed by your code snippet.\n\nExample:\n`/carbon fmt.Println(\"Hello\")`",
		"carbon_wait": "Generating image... please wait.",
		"carbon_error": "Failed to generate image, check your code.",
		"back_btn": "Back",
	},
}

func GetLang(userID int64) string {
	lang, err := RDB.Get(context.Background(), fmt.Sprintf("lang:%d", userID)).Result()
	if err != nil {
		return "ar"
	}
	return lang
}
