import 'package:flutter/widgets.dart';
import 'package:retrovibed/timex.dart' as timex;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/inputs.dart' as inputs;
import 'package:retrovibed/design.kit/typography/timerange.dart' as tr;
import 'ast.dart';
import 'parser.dart' as lucene_parser;
import 'parser.results.dart';

// A typed field definition.
//
// Binds a lucene field name to a setter callback.
// Drives autocomplete suggestions and value parsing — independent of any UI.
abstract class Field<T> {
  final String name;
  final T current;
  final T defaultValue;
  final void Function(T) setter;
  final Widget help;

  const Field(this.name, this.current, this.defaultValue, this.setter, {this.help = ds.HelpScope.None});

  bool get available => current == defaultValue;

  Field<T> withCurrent(T value);

  /// Returns a Field<T> with the
  ParserResult of(T value);

  /// Returns a Node representation of the provided value for this field.
  Node from(T value);

  // Parse a raw term value string (the part after the colon) into T.
  // Never fails — returns defaultValue when the raw string is unrecognised.
  T parse(String raw);

  // Autocomplete suggestions for the partial string typed after the colon.
  List<Suggestion> suggestions(String partial);

  // widget shown at the top of the autocomplete dropdown.
  // Return null to show only text suggestions.
  Widget render(T value, void Function(T) onChanged) => ds.Empty;

  // When non-null, the field completes immediately at ':' with this value,
  // skipping Value state entirely.
  T? get autocomplete => null;

  // Walk an AST node and call setter if this field owns it.
  // Returns true when the node was consumed.
  bool apply(Node node) {
    return switch (node) {
      Term t when t.field == name => _applyTerm(t.value),
      _ => false,
    };
  }

  bool _applyTerm(String raw) {
    setter(parse(raw));
    return true;
  }
}

// A field that matches a lucene Range node and dispatches min/max to two
// sub-fields of the same scalar type.
class RangeField<T> extends Field<({T min, T max})> {
  final Field<T> minField;
  final Field<T> maxField;
  final Widget Function(({T min, T max}))? display;

  const RangeField(
    String name,
    ({T min, T max}) current,
    ({T min, T max}) defaultValue,
    void Function(({T min, T max})) setter,
    this.minField,
    this.maxField, {
    this.display,
    Widget help = ds.HelpScope.None,
  }) : super(name, current, defaultValue, setter, help: help);

  @override
  RangeField<T> withCurrent(({T min, T max}) value) => RangeField(
    name,
    value,
    defaultValue,
    setter,
    minField.withCurrent(value.min),
    maxField.withCurrent(value.max),
    display: display,
    help: help,
  );

  @override
  Node from(({T min, T max}) value) => Range(
    name,
    min: minField.from(value.min) is Term ? (minField.from(value.min) as Term).value : null,
    max: maxField.from(value.max) is Term ? (maxField.from(value.max) as Term).value : null,
  );

  @override
  ({T min, T max}) parse(String raw) {
    final q = lucene_parser.parse('f:$raw');
    if (q.clauses.isEmpty) return defaultValue;
    final node = q.clauses.first.node;
    if (node is! Range) return defaultValue;
    return (
      min: node.min != null ? minField.parse(node.min!) : current.min,
      max: node.max != null ? maxField.parse(node.max!) : current.max,
    );
  }

  @override
  List<Suggestion> suggestions(String partial) => const [];

  @override
  ParserResult of(({T min, T max}) value) => ParserResultRange(withCurrent(value));

  @override
  bool apply(Node node) {
    if (node case Range r when r.field == name) {
      final min = r.min != null ? minField.parse(r.min!) : current.min;
      final max = r.max != null ? maxField.parse(r.max!) : current.max;
      setter((min: min, max: max));
      return true;
    }
    return false;
  }
}

// A parsed autocomplete suggestion.
class Suggestion {
  final Field field;
  final String label;
  final String completion; // full lucene text to insert, e.g. "date:[NOW-30d TO *]"
  final String? description;

  const Suggestion({
    required this.field,
    required this.label,
    required this.completion,
    this.description,
  });
}

// Pure autocomplete function — stateless, driven by the field list.
List<Suggestion> complete(String partial, List<Field> fields) {
  final colon = partial.indexOf(':');

  if (colon == -1) {
    final prefix = partial.toLowerCase();
    return [
      for (final f in fields)
        if (f.name.startsWith(prefix)) Suggestion(field: f, label: f.name, completion: '${f.name}:'),
    ];
  }

  final fieldName = partial.substring(0, colon);
  final valuePartial = partial.substring(colon + 1);
  final field = fields.where((f) => f.name == fieldName).firstOrNull;
  return field?.suggestions(valuePartial) ?? [];
}

// Apply a fully-parsed Query against a list of field definitions.
List<String> applyQuery(Query query, List<Field> fields) {
  final unmatched = <String>[];
  for (final clause in query.clauses) {
    final consumed = fields.any((f) => f.apply(clause.node));
    if (!consumed) {
      switch (clause.node) {
        case Term t when t.field.isEmpty:
          unmatched.add(t.value);
        default:
          break;
      }
    }
  }
  return unmatched;
}

class Mode extends Field<bool> {
  const Mode(
    String name,
    bool current,
    bool defaultValue,
    void Function(bool) setter, {
    Widget help = ds.HelpScope.None,
  }) : super(name, current, defaultValue, setter, help: help);

  factory Mode.auto(
    String name,
    bool defaultValue,
    void Function(bool) setter, {
    Widget help = ds.HelpScope.None,
  }) {
    return Mode(name, defaultValue, defaultValue, setter, help: help);
  }

  @override
  bool? get autocomplete => !defaultValue;

  @override
  Mode withCurrent(bool value) => Mode(name, value, defaultValue, setter, help: help);

  @override
  Node from(bool value) => Term(name, value.toString());

  @override
  bool parse(String raw) => switch (raw.toLowerCase()) {
    'true' || '1' || 'yes' => true,
    'false' || '0' || 'no' => false,
    _ => defaultValue,
  };

  @override
  List<Suggestion> suggestions(String partial) => [
    Suggestion(field: this, label: 'on', completion: '$name:true'),
    Suggestion(field: this, label: 'off', completion: '$name:false'),
  ].where((s) => s.label.startsWith(partial)).toList();

  @override
  ParserResult of(bool value) {
    return ParserResultMode(Mode(name, value, defaultValue, setter, help: help));
  }
}

class Boolean extends Field<bool> {
  const Boolean(
    String name,
    bool current,
    bool defaultValue,
    void Function(bool) setter, {
    Widget help = ds.HelpScope.None,
  }) : super(name, current, defaultValue, setter, help: help);

  factory Boolean.auto(
    String name,
    bool defaultValue,
    void Function(bool) setter, {
    Widget help = ds.HelpScope.None,
  }) {
    return Boolean(name, defaultValue, defaultValue, setter, help: help);
  }

  @override
  bool? get autocomplete => !defaultValue;

  @override
  Boolean withCurrent(bool value) => Boolean(name, value, defaultValue, setter, help: help);

  @override
  Node from(bool value) => Term(name, value.toString());

  @override
  bool parse(String raw) => switch (raw.toLowerCase()) {
    'true' || '1' || 'yes' => true,
    'false' || '0' || 'no' => false,
    _ => defaultValue,
  };

  @override
  List<Suggestion> suggestions(String partial) => [
    Suggestion(field: this, label: 'on', completion: '$name:true'),
    Suggestion(field: this, label: 'off', completion: '$name:false'),
  ].where((s) => s.label.startsWith(partial)).toList();

  @override
  ParserResult of(bool value) {
    return ParserResultBool(Boolean(name, value, defaultValue, setter, help: help));
  }
}

class Timestamp extends Field<DateTime> {
  const Timestamp(
    String name,
    DateTime current,
    DateTime defaultValue,
    void Function(DateTime) setter, {
    Widget help = ds.HelpScope.None,
  }) : super(name, current, defaultValue, setter, help: help);

  factory Timestamp.auto(
    String name,
    DateTime defaultValue,
    void Function(DateTime) setter, {
    Widget help = ds.HelpScope.None,
  }) {
    return Timestamp(name, defaultValue, defaultValue, setter, help: help);
  }

  @override
  Timestamp withCurrent(DateTime value) => Timestamp(name, value, defaultValue, setter, help: help);

  @override
  Node from(DateTime value) => Term(name, value.toIso8601String());

  @override
  DateTime parse(String raw) => timex.iso8601(raw, empty: defaultValue);

  @override
  List<Suggestion> suggestions(String partial) => const [];

  @override
  Widget render(DateTime value, void Function(DateTime) onChanged) =>
      inputs.DateInput(value: value, onChanged: onChanged);

  @override
  ParserResult of(DateTime value) => ParserResultTimestamp(withCurrent(value));
}

class DateRange extends RangeField<DateTime> {
  static List<({String label, timex.Range range})> get _presets {
    final now = timex.now();
    return [
      (
        label: 'last 7 days',
        range: timex.Range(now.subtract(const Duration(days: 7)), now),
      ),
      (
        label: 'last 30 days',
        range: timex.Range(now.subtract(const Duration(days: 30)), now),
      ),
      (
        label: 'last year',
        range: timex.Range(now.subtract(const Duration(days: 365)), now),
      ),
    ];
  }

  DateRange(
    String name,
    timex.Range v,
    timex.Range defaultValue,
    void Function(timex.Range) setter, {
    Widget help = ds.HelpScope.None,
  }) : super(
        name,
        (min: v.begin, max: v.end),
        (min: defaultValue.begin, max: defaultValue.end),
        (p) => setter(timex.Range(p.min, p.max)),
        Timestamp(
          name,
          v.begin,
          defaultValue.begin,
          (d) => setter(timex.Range(d, v.end)),
        ),
        Timestamp(
          name,
          v.end,
          defaultValue.end,
          (d) => setter(timex.Range(v.begin, d)),
        ),
        display: (v) => Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text('$name: '),
            tr.TimeRange(timex.Range(v.min, v.max)),
          ],
        ),
        help: help,
      );

  factory DateRange.auto(
    String name,
    timex.Range defaultValue,
    void Function(timex.Range) setter, {
    Widget help = ds.HelpScope.None,
  }) {
    return DateRange(name, defaultValue, defaultValue, setter, help: help);
  }

  @override
  DateRange withCurrent(({DateTime min, DateTime max}) value) => DateRange(
    name,
    timex.Range(value.min, value.max),
    timex.Range(defaultValue.min, defaultValue.max),
    (r) => setter((min: r.begin, max: r.end)),
    help: help,
  );

  @override
  Node from(({DateTime min, DateTime max}) value) => Range(
    name,
    min: value.min.toIso8601String(),
    max: value.max == timex.inf ? null : value.max.toIso8601String(),
  );

  @override
  List<Suggestion> suggestions(String partial) => [
    for (final p in _presets)
      if (p.label.contains(partial))
        Suggestion(
          field: this,
          label: p.label,
          completion: '$name:[${p.range.begin.toIso8601String()} TO ${p.range.end.toIso8601String()}]',
        ),
  ];

  @override
  Widget render(
    ({DateTime min, DateTime max}) value,
    void Function(({DateTime min, DateTime max})) onChanged,
  ) {
    return inputs.DateRangeInput(
      value: timex.Range(value.min, value.max),
      onChanged: (r) => onChanged((min: r.begin, max: r.end)),
    );
  }
}

class Number extends Field<num> {
  const Number(
    String name,
    num v,
    num defaultValue,
    void Function(num) setter, {
    Widget help = ds.HelpScope.None,
  }) : super(name, v, defaultValue, setter, help: help);

  factory Number.auto(
    String name,
    num defaultValue,
    void Function(num) setter, {
    Widget help = ds.HelpScope.None,
  }) {
    return Number(name, defaultValue, defaultValue, setter, help: help);
  }

  @override
  Number withCurrent(num value) => Number(name, value, defaultValue, setter, help: help);

  @override
  Node from(num value) => Term(name, value.toString());

  @override
  num parse(String raw) => num.tryParse(raw.replaceFirst(RegExp(r'^[><=]+'), '')) ?? defaultValue;

  @override
  List<Suggestion> suggestions(String partial) => const [];

  @override
  ParserResult of(num value) {
    return ParserResultNumeric(Number(name, value, defaultValue, setter, help: help));
  }
}

class Elapsed extends Field<Duration> {
  static final _presets = [
    (label: '30 min', duration: const Duration(minutes: 30)),
    (label: '60 min', duration: const Duration(minutes: 60)),
    (label: '90 min', duration: const Duration(minutes: 90)),
    (label: '2 hours', duration: const Duration(hours: 2)),
  ];

  const Elapsed(
    String name,
    Duration v,
    Duration defaultValue,
    void Function(Duration) setter, {
    Widget help = ds.HelpScope.None,
  }) : super(name, v, defaultValue, setter, help: help);

  factory Elapsed.auto(
    String name,
    Duration defaultValue,
    void Function(Duration) setter, {
    Widget help = ds.HelpScope.None,
  }) {
    return Elapsed(name, defaultValue, defaultValue, setter, help: help);
  }

  @override
  Elapsed withCurrent(Duration value) => Elapsed(name, value, defaultValue, setter, help: help);

  @override
  Node from(Duration value) {
    return Term(name, timex.durations.iso8601(value));
  }

  @override
  Duration parse(String raw) => timex.durations.tryParse(raw) ?? defaultValue;

  @override
  List<Suggestion> suggestions(String partial) => [
    for (final p in _presets)
      if (p.label.contains(partial))
        Suggestion(
          field: this,
          label: p.label,
          completion: '$name:${timex.durations.iso8601(p.duration)}',
        ),
  ];

  @override
  ParserResult of(Duration value) {
    return ParserResultDuration(Elapsed(name, value, defaultValue, setter, help: help));
  }
}
