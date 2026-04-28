import 'package:flutter/material.dart';
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/storage.dart' as storage;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'rss.pb.dart';

class FeedNew extends StatelessWidget {
  final community = uuidx.random();
  final Feed current;
  final Function(Feed)? onChange;

  FeedNew({super.key, Feed? current, this.onChange}) : current = current ?? (Feed.create()..autodownload = false);

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    return forms.Container(
      Column(
        mainAxisSize: MainAxisSize.min,
        mainAxisAlignment: MainAxisAlignment.start,
        crossAxisAlignment: CrossAxisAlignment.stretch,
        spacing: defaults.spacing,
        children: [
          forms.Field(
            input: TextFormField(
              decoration: InputDecoration(
                helperText: current.url.isEmpty ? "description" : current.url,
              ),
              initialValue: current.description,
              onChanged: (v) => onChange?.call(current..description = v),
            ),
          ),
          forms.Field(
            input: TextFormField(
              initialValue: current.url,
              decoration: InputDecoration(helperText: "url"),
              onChanged: (v) => onChange?.call(current..url = v),
            ),
          ),
          storage.SeedTypography(
            current.encryptionSeed,
            classifier: storage.Classifier(
              community: community,
              personal: uuidx.max(),
            ),
            onChange: (v) => onChange?.call(current..encryptionSeed = v.id),
          ),
          ds.Container(
            Wrap(
              children: [
                forms.Checkbox(
                  const Text("autodownload"),
                  value: current.autodownload,
                  onChanged: (v) {
                    onChange?.call(
                      current..autodownload = (v ?? current.autodownload),
                    );
                  },
                ),
                forms.Checkbox(
                  const Text("autoarchive"),
                  value: current.autoarchive,
                  onChanged: (v) {
                    onChange?.call(
                      current..autoarchive = (v ?? current.autoarchive),
                    );
                  },
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
