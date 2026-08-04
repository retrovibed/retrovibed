import 'package:iso_duration_parser/iso_duration_parser.dart' as derp;

final DateTime inf = DateTime.fromMillisecondsSinceEpoch(253402300799999, isUtc: true).toUtc();
final DateTime epoch = DateTime.fromMillisecondsSinceEpoch(0, isUtc: true).toUtc();
final DateTime neginf = DateTime.utc(0, 1, 1, 1, 1, 1); // matches backend timex.RFC3339NegInf()

class durations {
  static Duration? tryParse(
    String input, {
    Duration? fallback = const Duration(),
  }) {
    var tmp = derp.IsoDuration.tryParse(input);
    var ts = now();
    return tmp?.withDate(ts).difference(ts) ?? fallback;
  }

  static String iso8601(Duration input) {
    final isoDuration = derp.IsoDuration(
      days: input.inDays.toDouble(),
      hours: input.inHours.remainder(24),
      minutes: input.inMinutes.remainder(60),
      seconds: input.inSeconds.remainder(60),
    );

    return isoDuration.toIso();
  }
}

class Range {
  final DateTime begin;
  final DateTime end;

  const Range(this.begin, this.end);

  factory Range.latest(Duration d) {
    final ts = now();
    return Range(ts.subtract(d), ts);
  }

  factory Range.everything() {
    return Range(neginf, now());
  }

  @override
  bool operator ==(Object other) =>
      identical(this, other) || other is Range && begin == other.begin && end == other.end;

  @override
  int get hashCode => Object.hash(begin, end);
}

DateTime iso8601(String ts, {DateTime? empty}) {
  return ts.isEmpty ? empty ?? epoch : DateTime.parse(ts).toUtc();
}

String formatISO8601(DateTime ts) {
  return ts.toIso8601String();
}

DateTime now({bool utc = true}) {
  final ts = DateTime.now();
  return utc ? ts.toUtc() : ts;
}

DateTime min(Iterable<DateTime> dates, {DateTime? fallback}) {
  if (dates.isEmpty) {
    return neginf;
  }

  DateTime m = inf;

  for (DateTime date in dates) {
    if (!date.isBefore(m)) {
      continue;
    }

    m = date;
  }

  return m;
}

DateTime max(Iterable<DateTime> dates, {DateTime? fallback}) {
  if (dates.isEmpty) {
    return inf;
  }

  DateTime m = neginf;

  for (DateTime date in dates) {
    if (!date.isAfter(m)) {
      continue;
    }

    m = date;
  }

  return m;
}
