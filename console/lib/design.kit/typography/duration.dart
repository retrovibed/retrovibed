import 'dart:core' as core;
import 'package:flutter/material.dart';

typedef DurationFormatter = Widget Function(core.Duration duration);

class Duration extends StatelessWidget {
  final core.Duration duration;
  final DurationFormatter formatter;

  const Duration(this.duration, {super.key, this.formatter = relative});

  factory Duration.until(
    core.DateTime timestamp, {
    DurationFormatter formatter = relative,
  }) {
    return Duration(
      timestamp.difference(core.DateTime.now()),
      formatter: formatter,
    );
  }

  factory Duration.untilISO8601(
    core.String isoTimestamp, {
    DurationFormatter formatter = relative,
  }) {
    return Duration.until(
      core.DateTime.parse(isoTimestamp),
      formatter: formatter,
    );
  }

  static core.String _padTwoDigits(core.int n) => n.toString().padLeft(2, '0');

  /// Formats duration as HH:MM:SS (e.g., "02:30:45" or "-01:15:00")
  static Widget elapsed(core.Duration duration) {
    core.String sign = '';
    core.Duration positiveDuration = duration;

    if (duration.isNegative) {
      sign = '-';
      positiveDuration = duration.abs();
    }

    final core.String hours = positiveDuration.inHours.toString().padLeft(
      2,
      '0',
    );
    final core.String minutes = _padTwoDigits(
      positiveDuration.inMinutes.remainder(60),
    );
    final core.String seconds = _padTwoDigits(
      positiveDuration.inSeconds.remainder(60),
    );

    return Text('$sign$hours:$minutes:$seconds');
  }

  /// Formats duration as human-readable relative time (e.g., "5m ago", "2h ago", "3d ago")
  static Widget ago(core.Duration duration) {
    final d = duration.abs();

    final core.String text;
    if (d.inDays > 0) {
      text = '${d.inDays}d ago';
    } else if (d.inHours > 0) {
      text = '${d.inHours}h ago';
    } else if (d.inMinutes > 0) {
      text = '${d.inMinutes}m ago';
    } else {
      text = 'just now';
    }

    return Text(text);
  }

  /// Formats duration as human-readable future time (e.g., "in 5m", "in 2h", "in 3d")
  static Widget future(core.Duration duration) {
    final d = duration.abs();

    final core.String text;
    if (d.inDays > 0) {
      text = 'in ${d.inDays}d';
    } else if (d.inHours > 0) {
      text = 'in ${d.inHours}h';
    } else if (d.inMinutes > 0) {
      text = 'in ${d.inMinutes}m';
    } else {
      text = 'now';
    }

    return Text(text);
  }

  /// Formats duration as "ago" for past or "in X" for future based on sign
  static Widget relative(core.Duration duration) {
    if (duration.isNegative) {
      return ago(duration);
    } else {
      return future(duration);
    }
  }

  @core.override
  Widget build(BuildContext context) {
    return formatter(duration);
  }
}
