import 'package:flutter/material.dart';
import 'package:retrovibed/design.kit/typography/duration.dart' as dur;
import 'package:retrovibed/design.kit/typography/timestamp.dart' as ts;
import './field.dart';
import './parser.states.dart' show Parser;

sealed class ParserResult extends StatelessWidget {
  const ParserResult({super.key});

  static const ParserResult close = _ParserResultClose();

  Widget? edit(void Function(ParserResult) onChanged);

  void reset(Parser parser);
  void apply(Parser parser);
}

class _ParserResultClose extends ParserResult {
  const _ParserResultClose();

  @override
  Widget build(BuildContext context) => const SizedBox.shrink();

  @override
  Widget? edit(void Function(ParserResult) onChanged) => null;

  @override
  void reset(Parser parser) {}

  @override
  void apply(Parser parser) {}
}

class ParserResultMode extends ParserResult {
  final Field<bool> field;
  const ParserResultMode(this.field, {super.key});

  @override
  Widget build(BuildContext context) {
    return Text(field.name);
  }

  @override
  Widget? edit(void Function(ParserResult) onChanged) {
    return null; // cannot be edited, only removed.
  }

  @override
  void reset(Parser parser) {
    field.apply(field.from(field.defaultValue));
    parser.replace(field.withCurrent(field.defaultValue));
  }

  @override
  void apply(Parser parser) {
    field.apply(field.from(field.current));
    parser.replace(field.withCurrent(field.current));
  }
}

class ParserResultBool extends ParserResult {
  final Field<bool> field;
  const ParserResultBool(this.field, {super.key});

  @override
  Widget build(BuildContext context) {
    return Text(field.name);
  }

  @override
  Widget? edit(void Function(ParserResult) onChanged) {
    return null; // cannot be edited, only removed.
  }

  @override
  void reset(Parser parser) {
    field.apply(field.from(field.defaultValue));
    parser.replace(field.withCurrent(field.defaultValue));
  }

  @override
  void apply(Parser parser) {
    field.apply(field.from(field.current));
    parser.replace(field.withCurrent(field.current));
  }
}

class ParserResultRange<T> extends ParserResult {
  final RangeField<T> field;
  const ParserResultRange(this.field, {super.key});

  @override
  Widget build(BuildContext context) {
    return field.display?.call(field.current) ?? Text('${field.name}: ${field.current}');
  }

  @override
  Widget? edit(void Function(ParserResult) onChanged) {
    return field.render(field.current, (v) => onChanged(field.of(v)));
  }

  @override
  void reset(Parser parser) {
    field.apply(field.from(field.defaultValue));
    parser.replace(field.withCurrent(field.defaultValue));
  }

  @override
  void apply(Parser parser) {
    field.apply(field.from(field.current));
    parser.replace(field.withCurrent(field.current));
  }
}

class ParserResultTimestamp extends ParserResult {
  final Field<DateTime> field;
  const ParserResultTimestamp(this.field, {super.key});

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [Text('${field.name}: '), ts.Timestamp(field.current)],
    );
  }

  @override
  Widget? edit(void Function(ParserResult) onChanged) => field.render(field.current, (v) => onChanged(field.of(v)));

  @override
  void reset(Parser parser) {
    field.apply(field.from(field.defaultValue));
    parser.replace(field.withCurrent(field.defaultValue));
  }

  @override
  void apply(Parser parser) {
    field.apply(field.from(field.current));
    parser.replace(field.withCurrent(field.current));
  }
}

class ParserResultNumeric extends ParserResult {
  final Field<num> field;
  const ParserResultNumeric(this.field, {super.key});

  @override
  Widget build(BuildContext context) {
    return Text('${field.name}: ${field.current}');
  }

  @override
  Widget? edit(void Function(ParserResult) onChanged) {
    return null; // cannot be edited, only removed.
  }

  @override
  void reset(Parser parser) {
    field.apply(field.from(field.defaultValue));
    parser.replace(field.withCurrent(field.defaultValue));
  }

  @override
  void apply(Parser parser) {
    field.apply(field.from(field.current));
    parser.replace(field.withCurrent(field.current));
  }
}

class ParserResultDuration extends ParserResult {
  final Field<Duration> field;
  const ParserResultDuration(this.field, {super.key});

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Text('${field.name}: '),
        dur.Duration(field.current, formatter: dur.Duration.elapsed),
      ],
    );
  }

  @override
  Widget? edit(void Function(ParserResult) onChanged) {
    return null; // cannot be edited, only removed.
  }

  @override
  void reset(Parser parser) {
    field.apply(field.from(field.defaultValue));
    parser.replace(field.withCurrent(field.defaultValue));
  }

  @override
  void apply(Parser parser) {
    field.apply(field.from(field.current));
    parser.replace(field.withCurrent(field.current));
  }
}
