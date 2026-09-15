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
    this.autofocus = true,
  }) : firstDate = firstDate ?? timex.neginf,
       lastDate = lastDate ?? timex.inf;

  @override
  State<DateRangeInput> createState() => _DateRangeInputState(value);
}

class _DateRangeInputState extends State<DateRangeInput> {
  Widget _picker = ds.Empty;
  DateTime _current = timex.epoch;
  timex.Range _pending;

  // The range last handed to widget.onChanged (via _apply, e.g. pressing
  // Enter). Since widget.value doesn't update until the caller rebuilds this
  // widget with the new value — and a caller that treats onChanged as "this
  // filter is now committed" may instead tear this widget down entirely —
  // deactivate() must not compare _pending against widget.value alone, or it
  // re-fires onChanged with a range it already reported, racing the caller's
  // own in-flight handling of the first call.
  timex.Range? _applied;

  // parentScope (rather than the default closedLoop) lets Tab/Shift+Tab
  // escape to whatever's next outside this widget — e.g. a field's preset
  // suggestion list shown alongside the picker — instead of endlessly
  // cycling begin/end/calendar in place.
  final FocusScopeNode _focusScopeNode = FocusScopeNode(
    traversalEdgeBehavior: TraversalEdgeBehavior.parentScope,
  );

  _DateRangeInputState(this._pending);

  @override
  void initState() {
    super.initState();
    postframe(() => _showBegin());
  }

  @override
  void dispose() {
    _focusScopeNode.dispose();
    super.dispose();
  }

  @override
  void deactivate() {
    if (_pending != widget.value && _pending != _applied) {
      postframe(() => widget.onChanged(_pending));
    }
    super.deactivate();
  }

  void _apply() {
    widget.onChanged(_pending);
    _applied = _pending;
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
        initialDate: timex.max([_pending.begin.toLocal(), timex.now()]),
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
      node: _focusScopeNode,
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
                  onPressed: _showBegin,
                  style: timex.neginf != _current && _current == _pending.begin ? activestyle : null,
                  child: typography.Timestamp(_pending.begin, neginf: Text("")),
                ),
              ),
              const Text('–'),
              Expanded(
                child: TextButton(
                  onPressed: _showEnd,
                  style: timex.inf != _current && _current == _pending.end ? activestyle : null,
                  child: typography.Timestamp(_pending.end),
                ),
              ),
            ],
          ),
          Focus(autofocus: widget.autofocus, child: _picker),
        ],
      ),
    );
  }
}
