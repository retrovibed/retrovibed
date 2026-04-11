import 'package:flutter/material.dart';
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/media.dart' as media;

class RecentEdit extends StatelessWidget {
  final media.RecentRecordRequest current;
  final EdgeInsets? padding;
  final EdgeInsets? margin;
  final BoxDecoration? decoration;
  final BoxConstraints? constraints;
  final Alignment? alignment;
  final Color? background;
  final Clip clipBehavior;

  const RecentEdit(
    this.current, {
    super.key,
    this.padding,
    this.margin,
    this.decoration,
    this.constraints,
    this.alignment,
    this.background,
    this.clipBehavior = Clip.none,
  });

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final pos = current.position.toInt();
    final dur = current.duration.toInt();
    final progress = dur > 0 ? (pos / dur).clamp(0.0, 1.0) : 0.0;

    return ds.Container(
      padding: padding ?? defaults.padding,
      margin: margin,
      decoration: decoration,
      constraints: constraints,
      alignment: alignment,
      background: background,
      clipBehavior: clipBehavior,
      Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          forms.Field(label: Text('id'), input: Text(current.id)),
          forms.Field(
            label: Text('description'),
            input: Text(current.media.description),
          ),
          forms.Field(
            label: Text('mimetype'),
            input: Text(current.media.mimetype),
          ),
          forms.Field(
            label: Text('position'),
            input: ds.Duration(
              Duration(milliseconds: pos),
              formatter: ds.Duration.elapsed,
            ),
          ),
          forms.Field(
            label: Text('duration'),
            input: ds.Duration(
              Duration(milliseconds: dur),
              formatter: ds.Duration.elapsed,
            ),
          ),
          forms.Field(
            label: Text('progress'),
            input: LinearProgressIndicator(value: progress),
          ),
          forms.Field(
            label: Text('query'),
            input: Text(current.query.query),
          ),
          forms.Field(
            label: Text('mimetypes'),
            input: Text(current.query.mimetypes.join(', ')),
          ),
          forms.Field(
            label: Text('limit'),
            input: Text(current.query.limit.toString()),
          ),
          forms.Field(
            label: Text('offset'),
            input: Text(current.query.offset.toString()),
          ),
          forms.Field(
            label: Text('adult'),
            input: Text(current.query.adult.toString()),
          ),
          forms.Field(
            label: Text('hidden'),
            input: Text(current.query.hidden.toString()),
          ),
        ],
      ),
    );
  }
}
