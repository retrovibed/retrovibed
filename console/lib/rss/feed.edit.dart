import 'package:flutter/material.dart';
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/storage.dart' as storage;
import 'package:retrovibed/uuidx.dart' as uuidx;
import './rss.pb.dart';

class Edit extends StatelessWidget {
  final Feed current;
  final Function(Feed)? onChange;
  final EdgeInsets? padding;
  Edit({super.key, Feed? current, this.onChange, this.padding})
    : current = current ?? (Feed.create()..autodownload = false);

  @override
  Widget build(BuildContext context) {
    return forms.Container(
      Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          TextFormField(
            decoration: InputDecoration(helperText: "description"),
            initialValue: current.description,
            autofocus: true,
            onChanged: (v) => onChange?.call(current..description = v),
          ),
          TextFormField(
            decoration: InputDecoration(helperText: "url"),
            initialValue: current.url,
            onChanged: (v) => onChange?.call(current..url = v),
          ),
          storage.SeedTypography(
            current.encryptionSeed,
            classifier: storage.Classifier(
              community: current.encryptionSeed,
              personal: uuidx.max(),
            ),
            onChange: (v) => onChange?.call(current..encryptionSeed = v.id),
          ),
          Wrap(
            children: [
              forms.Checkbox(
                Text("autodownload"),
                alignment: Alignment.topLeft,
                value: current.autodownload,
                onChanged: (v) {
                  onChange?.call(
                    current..autodownload = (v ?? current.autodownload),
                  );
                },
              ),
              forms.Checkbox(
                Text("autoarchive"),
                alignment: Alignment.centerLeft,
                value: current.autoarchive,
                onChanged: (v) {
                  onChange?.call(
                    current..autoarchive = (v ?? current.autoarchive),
                  );
                },
              ),
              Tooltip(
                message:
                    "support this open source community by providing distribution when autodownload is enabled, and financially when autoarchive is enabled",
                child: forms.Checkbox(
                  Text("contribution"),
                  alignment: Alignment.centerLeft,
                  value: current.contributing,
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}
