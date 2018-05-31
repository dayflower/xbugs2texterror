
target:
	curl https://raw.githubusercontent.com/spotbugs/spotbugs/release-3.1/spotbugs/etc/messages.xml | python message2code.py -n BugDescriptionEn > ../messages_en.go
	curl https://raw.githubusercontent.com/spotbugs/spotbugs/release-3.1/spotbugs/etc/messages_ja.xml | python message2code.py -n BugDescriptionJa > ../messages_ja.go
	curl https://raw.githubusercontent.com/spotbugs/spotbugs/release-3.1/spotbugs/etc/messages_fr.xml | python message2code.py -n BugDescriptionFr > ../messages_fr.go
	( cd .. && go fmt )

