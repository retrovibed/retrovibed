import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import './../../design.kit/theme.defaults.dart';

// Unit choices and their Duration multipliers.
const _units = [
  (label: 's', duration: Duration(seconds: 1)),
  (label: 'm', duration: Duration(minutes: 1)),
  (label: 'hr', duration: Duration(hours: 1)),
];

class DurationInput extends StatefulWidget {
  final Duration value;
  final ValueChanged<Duration> onChanged;

  const DurationInput({
    super.key,
    required this.value,
    required this.onChanged,
  });

  @override
  State<DurationInput> createState() => _DurationInputState();
}

class _DurationInputState extends State<DurationInput> {
  late ({String label, Duration duration}) _unit;

  @override
  void initState() {
    super.initState();
    _unit = _unitFor(widget.value);
  }

  // Pick the largest unit that divides the duration evenly, defaulting to s.
  ({String label, Duration duration}) _unitFor(Duration d) {
    if (d.inMinutes > 0 && d.inMinutes % 60 == 0) return _units[2]; // hr
    if (d.inSeconds > 0 && d.inSeconds % 60 == 0) return _units[1]; // m
    return _units[0]; // s
  }

  String _initialText() {
    final amount = widget.value.inSeconds ~/ _unit.duration.inSeconds;
    return amount > 0 ? amount.toString() : '';
  }

  void _emit(String text) {
    final n = int.tryParse(text);
    if (n == null || n <= 0) return;
    widget.onChanged(_unit.duration * n);
  }

  @override
  Widget build(BuildContext context) {
    final defaults = Defaults.of(context);

    return Row(
      mainAxisSize: MainAxisSize.min,
      spacing: defaults.spacing,
      children: [
        Flexible(
          child: TextFormField(
            key: ValueKey(_unit),
            initialValue: _initialText(),
            keyboardType: TextInputType.number,
            inputFormatters: [FilteringTextInputFormatter.digitsOnly],
            decoration: InputDecoration(isDense: true, hintText: _unit.label),
            onFieldSubmitted: _emit,
          ),
        ),
        DropdownButton<({String label, Duration duration})>(
          value: _unit,
          underline: const SizedBox(),
          onChanged: (unit) {
            if (unit == null) return;
            setState(() => _unit = unit);
          },
          items: _units.map((u) => DropdownMenuItem(value: u, child: Text(u.label))).toList(),
        ),
      ],
    );
  }
}
