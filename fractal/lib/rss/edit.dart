import 'package:flutter/material.dart';
import 'package:fractal/design.kit/forms.dart' as forms;
import './rss.pb.dart';

class Edit extends StatefulWidget {
  final Feed feed;
  final Function(Feed)? onChange;

  Edit({super.key, Feed? feed, this.onChange})
    : feed =
          feed ?? Feed.create()
            ..autodownload = false;

  @override
  State<Edit> createState() => EditView();
}

class EditView extends State<Edit> {
  // Feed feed;

  // EditView(this.feed);

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
              onChanged:
                  // (v) => setState(() {
                  //   widget.feed.description = v;
                  // }),
                  (v) => widget.onChange?.call(widget.feed..description = v),
            ),
          ),
          forms.Field(
            label: Text("url"),
            input: TextFormField(
              initialValue: widget.feed.url,
              onChanged:
                  // (v) => setState(() {
                  //   widget.feed.url = v;
                  // }),
                  (v) => widget.onChange?.call(widget.feed..url = v),
            ),
          ),
          forms.Field(
            label: Text("autodownload"),
            input: Row(
              children: [
                Checkbox(
                  value: widget.feed.autodownload,
                  onChanged:
                      (v) => widget.onChange?.call(
                        widget.feed
                          ..autodownload = v ?? widget.feed.autodownload,
                      ),
                  // (v) => setState(
                  //   () =>
                  //       widget.feed.autodownload =
                  //           v ?? widget.feed.autodownload,
                  // ),
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
