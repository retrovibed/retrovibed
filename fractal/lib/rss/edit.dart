import 'package:flutter/material.dart';
import 'package:fractal/design.kit/forms.dart' as forms;
import './rss.pb.dart';

class Edit extends StatelessWidget {
  final Feed feed;
  final Function(Feed)? onChange;
  Edit({super.key, Feed? feed, this.onChange})
    : feed = feed ?? (Feed.create()..autodownload = false);

  @override
  Widget build(BuildContext context) {
    final theming = Theme.of(context);

    return Container(
      color: theming.scaffoldBackgroundColor,
      child: ConstrainedBox(
        constraints: BoxConstraints(minHeight: 128, minWidth: 128),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            forms.Field(
              label: Text("description"),
              input: TextFormField(
                initialValue: feed.description,
                onChanged: (v) => onChange?.call(feed..description = v),
              ),
            ),
            forms.Field(
              label: Text("url"),
              input: TextFormField(
                initialValue: feed.url,
                onChanged: (v) => onChange?.call(feed..url = v),
              ),
            ),
            forms.Field(
              label: Text("autodownload"),
              input: Row(
                children: [
                  Checkbox(
                    value: feed.autodownload,
                    onChanged:
                        (v) => onChange?.call(
                          feed..autodownload = (v ?? feed.autodownload),
                        ),
                  ),
                  Spacer(),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
