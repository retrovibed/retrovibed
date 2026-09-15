import 'package:flutter/material.dart';
import 'package:retrovibed/timex.dart' as timex;

class DateInput extends StatefulWidget {
  final DateTime? value;
  final ValueChanged<DateTime> onChanged;
  final DateTime? firstDate;
  final DateTime? lastDate;
  final bool autofocus;

  const DateInput({
    super.key,
    required this.onChanged,
    this.value,
    this.firstDate,
    this.lastDate,
    this.autofocus = true,
  });

  @override
  State<DateInput> createState() => _DateInputState();
}

class _DateInputState extends State<DateInput> {
  late DateTime _pending;

  @override
  void initState() {
    super.initState();
    _pending = (widget.value ?? timex.now()).toLocal();
  }

  @override
  Widget build(BuildContext context) {
    return Focus(
      autofocus: widget.autofocus,
      child: CalendarDatePicker(
        initialDate: _pending,
        firstDate: (widget.firstDate ?? timex.epoch).toLocal(),
        lastDate: (widget.lastDate ?? timex.inf).toLocal(),
        onDateChanged: (picked) {
          setState(() => _pending = picked);
          widget.onChanged(picked.toUtc());
        },
      ),
    );
  }
}
