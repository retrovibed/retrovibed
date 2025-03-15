import 'package:flutter/material.dart';
import 'package:fractal/design.kit/forms.dart' as forms;
import './rss.pb.dart';

class FeedNew extends StatelessWidget {
  final Feed current;
  final Function(Feed)? onChange;
  final EdgeInsetsGeometry? padding;
  FeedNew({super.key, Feed? current, this.onChange, this.padding})
    : current = current ?? (Feed.create()..autodownload = false);

  @override
  Widget build(BuildContext context) {
    final theming = Theme.of(context);

    return Container(
      padding: padding,
      color: theming.scaffoldBackgroundColor,
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          forms.Field(
            label: SelectableText("url"),
            input: TextFormField(
              initialValue: current.url,
              onChanged: (v) => onChange?.call(current..url = v),
            ),
          ),
          forms.Field(
            label: SelectableText("autodownload"),
            input: Row(
              children: [
                Checkbox(
                  value: current.autodownload,
                  onChanged: (v) {
                    onChange?.call(
                      current..autodownload = (v ?? current.autodownload),
                    );
                  },
                ),
                Spacer(),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
