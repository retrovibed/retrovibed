import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:retrovibed/timex.dart' as timex;

class Timestamp extends StatelessWidget {
  final DateTime timestamp;
  final String Function(DateTime)? format;
  final Widget leading;
  final Widget trailing;
  final Widget _inf = const Text("never", overflow: TextOverflow.ellipsis);
  final Widget _neginf = const Text("always", overflow: TextOverflow.ellipsis);

  const Timestamp(
    this.timestamp, {
    super.key,
    this.format,
    this.leading = const SizedBox(),
    this.trailing = const SizedBox(),
  });

  factory Timestamp.iso8601(
    String ts, {
    DateTime? empty,
    Widget leading = const SizedBox(),
    Widget trailing = const SizedBox(),
  }) {
    return Timestamp(
      timex.iso8601(ts, empty: empty),
      leading: leading,
      trailing: trailing,
    );
  }

  static String Function(DateTime) _defaultFormat(double width) {
    if (width < 120) {
      return DateFormat("M/d/y").format;
    } else if (width < 200) {
      return DateFormat("MMM d, y").format;
    } else if (width < 300) {
      return DateFormat("MMM d, y hh:mm a").format;
    }
    return DateFormat("y MMMM d EEEE hh:mm a").format;
  }

  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(
      builder: (context, constraints) {
        Widget content;
        if (timex.inf.difference(timestamp).inMilliseconds == 0) {
          content = _inf;
        } else if (timex.neginf.difference(timestamp).inMilliseconds == 0) {
          content = _neginf;
        } else {
          content = Text(
            (this.format ?? _defaultFormat(constraints.maxWidth))(
              timestamp.toLocal(),
            ),
            overflow: TextOverflow.ellipsis,
          );
        }

        return Row(
          mainAxisSize: MainAxisSize.min,
          children: [leading, Flexible(child: content), trailing],
        );
      },
    );
  }
}
