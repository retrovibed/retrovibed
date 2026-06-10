import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:retrovibed/designkit.dart' as ds;
import '../flutterx.dart';
import 'package:retrovibed/timex.dart' as timex;
import 'package:retrovibed/design.kit/typography.dart' as typography;

class DateRangeInput extends StatefulWidget {
  final timex.Range value;
  final ValueChanged<timex.Range> onChanged;
  final DateTime firstDate;
  final DateTime lastDate;
  final bool autofocus;

  DateRangeInput({
    super.key,
    required this.value,
    required this.onChanged,
    DateTime? firstDate,
    DateTime? lastDate,
    this.autofocus = false,
  }) : firstDate = firstDate ?? timex.neginf,
       lastDate = lastDate ?? timex.inf;

  @override
  State<DateRangeInput> createState() => _DateRangeInputState(value);
}

class _DateRangeInputState extends State<DateRangeInput> {
  Widget _picker = ds.Empty;
  DateTime _current = timex.epoch;
  timex.Range _pending;

  _DateRangeInputState(this._pending);

  @override
  void initState() {
    super.initState();
    postframe(() => _showBegin());
  }

  @override
  void deactivate() {
    if (_pending != widget.value) {
      postframe(() => widget.onChanged(_pending));
    }
    super.deactivate();
  }

  void _apply() {
    widget.onChanged(_pending);
    setState(() {
      _picker = ds.Empty;
      _current = timex.epoch;
    });
  }

  void _showBegin() {
    final firstDate = timex.min([widget.firstDate, _pending.begin]);
    final lastDate = timex.max([widget.lastDate, _pending.end]);

    setState(() {
      _current = _pending.begin;
      _picker = CalendarDatePicker(
        key: const ValueKey('begin'),
        initialDate: _pending.begin.toLocal(),
        firstDate: firstDate.toLocal(),
        lastDate: lastDate.toLocal(),
        onDateChanged: (d) {
          setState(() {
            _pending = timex.Range(d.toUtc(), _pending.end);
          });
        },
      );
    });
  }

  void _showEnd() {
    final firstDate = timex.min([widget.firstDate, _pending.begin]);
    final lastDate = timex.max([widget.lastDate, _pending.end]);
    setState(() {
      _current = _pending.end;
      _picker = CalendarDatePicker(
        key: const ValueKey('end'),
        initialDate: _pending.end.toLocal(),
        firstDate: firstDate.toLocal(),
        lastDate: lastDate.toLocal(),
        onDateChanged: (d) {
          setState(() {
            _pending = timex.Range(_pending.begin, d.toUtc());
          });
        },
      );
    });
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final theme = Theme.of(context);
    final activestyle = TextButton.styleFrom(
      backgroundColor: theme.colorScheme.primaryContainer,
    );

    return FocusScope(
      onFocusChange: (hasFocus) {
        if (!hasFocus && _pending != widget.value) {
          _apply();
        }
      },
      onKeyEvent: (node, event) {
        if (event is KeyDownEvent && event.logicalKey == LogicalKeyboardKey.enter && _pending != widget.value) {
          _apply();
          return KeyEventResult.handled;
        }
        return KeyEventResult.ignored;
      },
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Row(
            spacing: defaults.spacing,
            children: [
              Expanded(
                child: TextButton(
                  autofocus: widget.autofocus,
                  onPressed: _showBegin,
                  style: _current == _pending.begin ? activestyle : null,
                  child: typography.Timestamp(_pending.begin),
                ),
              ),
              const Text('–'),
              Expanded(
                child: TextButton(
                  onPressed: _showEnd,
                  style: _current == _pending.end ? activestyle : null,
                  child: typography.Timestamp(_pending.end),
                ),
              ),
            ],
          ),
          _picker,
        ],
      ),
    );
  }
}
