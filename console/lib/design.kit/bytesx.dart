class bytesx {
  final int value;
  static const int KiB = 1024;
  static const int MiB = 1048576;
  static const int GiB = 1073741824;
  static const int TiB = 1099511627776;

  const bytesx(this.value);

  static const Map<int, String> _valueToStringMap = {
    KiB: 'KiB',
    MiB: 'MiB',
    GiB: 'GiB',
    TiB: 'TiB',
  };

  static String getName(int value) {
    return _valueToStringMap[value] ?? "";
  }
}

extension ByteFormatter on bytesx {
  double toIEC600272() {
    // IEC600272 constants
    const int KiB = 1024;
    const int MiB = 1024 * KiB;
    const int GiB = 1024 * MiB;
    const int TiB = 1024 * GiB;
    const int PiB = 1024 * TiB;
    const int EiB = 1024 * PiB;

    int div = 1;

    if (value >= EiB) {
      div = EiB;
    } else if (value >= PiB) {
      div = PiB;
    } else if (value >= TiB) {
      div = TiB;
    } else if (value >= GiB) {
      div = GiB;
    } else if (value >= MiB) {
      div = MiB;
    } else if (value >= KiB) {
      div = KiB;
    }

    double convertedValue = (value / div).toDouble();
    return convertedValue;
  }

  String toIEC600272Format() {
    // IEC600272 constants
    const int KiB = 1024;
    const int MiB = 1024 * KiB;
    const int GiB = 1024 * MiB;
    const int TiB = 1024 * GiB;
    const int PiB = 1024 * TiB;
    const int EiB = 1024 * PiB;
    
    int div = 1;
    String suffix = "B";

    if (value >= EiB) {
      div = EiB;
      suffix = "EiB";
    } else if (value >= PiB) {
      div = PiB;
      suffix = "PiB";
    } else if (value >= TiB) {
      div = TiB;
      suffix = "TiB";
    } else if (value >= GiB) {
      div = GiB;
      suffix = "GiB";
    } else if (value >= MiB) {
      div = MiB;
      suffix = "MiB";
    } else if (value >= KiB) {
      div = KiB;
      suffix = "KiB";
    }

    double convertedValue = (value / div).toDouble();
    if (div != 1) {
      convertedValue = convertedValue.roundToDouble();
    }
    
    if (convertedValue == convertedValue.toInt().toDouble()) {
      return "${convertedValue.toInt()} $suffix";
    } else {
      return "${convertedValue.toStringAsFixed(1)} $suffix";
    }
  }

  String toSIFormat() {
    // SI constants
    const int KB_SI = 1000;
    const int MB_SI = 1000 * KB_SI;
    const int GB_SI = 1000 * MB_SI;
    const int TB_SI = 1000 * GB_SI;
    const int PB_SI = 1000 * TB_SI;
    const int EB_SI = 1000 * PB_SI;
    
    int div = 1;
    String suffix = "B";

    if (value >= EB_SI) {
      div = EB_SI;
      suffix = "EB";
    } else if (value >= PB_SI) {
      div = PB_SI;
      suffix = "PB";
    } else if (value >= TB_SI) {
      div = TB_SI;
      suffix = "TB";
    } else if (value >= GB_SI) {
      div = GB_SI;
      suffix = "GB";
    } else if (value >= MB_SI) {
      div = MB_SI;
      suffix = "MB";
    } else if (value >= KB_SI) {
      div = KB_SI;
      suffix = "kB";
    }

    double convertedValue = (value / div).toDouble();
    if (div != 1) {
      convertedValue = convertedValue.roundToDouble();
    }
    
    if (convertedValue == convertedValue.toInt().toDouble()) {
      return "${convertedValue.toInt()} $suffix";
    } else {
      return "${convertedValue.toStringAsFixed(1)} $suffix";
    }
  }
}