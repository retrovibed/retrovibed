import 'package:flutter/material.dart';
import 'package:fractal/design.kit/forms.dart' as forms;
import './rss.pb.dart';

class Edit extends StatefulWidget {
  final Feed feed;

  Edit({super.key, Feed? feed})
    : feed =
          feed ?? Feed.create()
            ..autodownload = true;

  @override
  State<Edit> createState() => EditView();
}

class EditView extends State<Edit> {
  @override
  Widget build(BuildContext context) {
    return ConstrainedBox(
      constraints: BoxConstraints(minHeight: 128, minWidth: 128),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          forms.Field(
            label: Text("description"),
            input: TextFormField(
              initialValue: widget.feed.description,
              onChanged: (v) => widget.feed.description = v,
            ),
          ),
          forms.Field(
            label: Text("url"),
            input: TextFormField(
              initialValue: widget.feed.url,
              onChanged: (v) => widget.feed.url = v,
            ),
          ),
          forms.Field(
            label: Text("autodownload"),
            input: Row(
              children: [
                Checkbox(
                  value: widget.feed.autodownload,
                  onChanged:
                      (v) => setState(
                        () =>
                            widget.feed.autodownload =
                                v ?? widget.feed.autodownload,
                      ),
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
