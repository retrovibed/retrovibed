import 'package:flutter/material.dart';
import 'package:language_code/language_code.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/forms.dart' as forms;

enum Bitrates {
  auto('Auto', 0),
  low('Low (64 kbps)', 64),
  medium('Medium (128 kbps)', 128),
  high('High (256 kbps)', 256),
  lossless('Lossless (FLAC/ALAC)', 1411);

  final String label;
  final int bitrate;

  const Bitrates(this.label, this.bitrate);
}

class SettingsAudio extends StatefulWidget {
  const SettingsAudio({super.key, this.margin = EdgeInsets.zero});
  final EdgeInsets margin;

  @override
  State<SettingsAudio> createState() => _SettingsAudioState();
}

class _SettingsAudioState extends State<SettingsAudio> {
  Bitrates _selectedBitrate = Bitrates.high;

  LanguageCodes _lang = LanguageCode.code;

  final List<LanguageCodes> _availableLanguages = LanguageCodes.values;

  List<DropdownMenuItem<T>> _buildEnumDropdownItems<T extends Enum>({
    required List<T> values,
    required String Function(T) label,
  }) {
    return values.map((item) {
      return DropdownMenuItem<T>(
        value: item,
        child: Text(label(item), overflow: TextOverflow.ellipsis),
      );
    }).toList();
  }

  List<DropdownMenuItem<LanguageCodes>> _buildLanguageDropdownItems() {
    return _availableLanguages.map((language) {
      return DropdownMenuItem<LanguageCodes>(
        value: language,
        child: Text(language.name, overflow: TextOverflow.ellipsis),
      );
    }).toList();
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);

    return ds.Container(
      alignment: Alignment.topLeft,
      padding: defaults.padding,
      margin: widget.margin,
      Column(
        mainAxisSize: MainAxisSize.min,
        spacing: defaults.spacing / 2,
        children: [
          forms.Field(
            label: const Text("Bitrate"),
            input: DropdownButton<Bitrates>(
              alignment: Alignment.topLeft,
              isExpanded: true,
              value: _selectedBitrate,
              items: _buildEnumDropdownItems<Bitrates>(
                values: Bitrates.values,
                label: (q) => q.label,
              ),
              onChanged: (Bitrates? v) {
                if (v == null) return;
                setState(() {
                  _selectedBitrate = v;
                });
              },
            ),
          ),
          forms.Field(
            label: const Text("Language"),
            input: DropdownButton<LanguageCodes>(
              alignment: Alignment.topLeft,
              isExpanded: true,
              value: _lang,
              items: _buildLanguageDropdownItems(),
              onChanged: (v) {
                if (v == null) return;
                setState(() {
                  _lang = v;
                });
              },
            ),
          ),
        ],
      ),
    );
  }
}
