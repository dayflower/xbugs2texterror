# xbugs2texterror

Output of findbugs/spotbugs translator for [reviewdog](https://github.com/haya14busa/reviewdog)

## Usage

    $ cat build/reports/spotbugs/main.xml \
      | xbugs2texterror -l=ja \
      | reviewdog -f=golint -name=spotbugs -ci=common

## Options

### Language

    -l=<language> (-lang=...) (default: en)

Specify language.
One of en/ja/fr.

## License

[MIT License](LICENSE.md)
