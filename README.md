# Monitoring UAM's Soil
This repository contains the implementation of the infrastructure behind the monitoring of
soil parameters undertaken as part of a research project at Universidad Autónoma de Madrid.

Put simply, this project:

1. Monitors the soil in areas near the campus, especially the soil's temperature, humidity and conductivity.
1. The information is sent over NB-IoT to a Lambda in AWS.
1. The Lambda parses and stores the information in an on-premise PostgreSQL instance.
1. The information is then displayed in a Grafana Cloud-based dashboard.

## Deployed Sensors
We have currently deployed the following sensors:

| **Serial Number** | **UART Interface's Password** |      **ICCID**      |
| :---------------: | :---------------------------: | :-----------------: |
|        XYZ        |            `307e11`           | 8988228066603107708 |
|  BC660GL25241383  |            `348411`           | 8988228066603107712 |
|  BC660GL25241385  |            `347811`           | 8988228066603107713 |
|  BC660GL25241386  |            `347612`           | 8988228066603107716 |


## Sensor Configuration
Sensors are configured over an UART through `picocom(1)`. Once can access their shell with

    $ picocom --echo --baud 9600 --parity n --stopbits 1 --databits 8 --echo /dev/the-char-device

The first input we must provide is the password as it appears in the table in the previous section.
Once that's done, we can apply the configuration we see fit. For our use case, this boils down
to:

    # Set the uplink interval to 15 minutes
    AT+TDC=900

    # Set the payload type to JSON/UDP
    AT+PRO=2,5

    # Point to 1nce's UDP broker
    AT+SERVADDR=udp.os.1nce.com,4445 / 81.45.14.171,5005

    # Configure 1nce's APN
    AT+APN=iot.1nce.net / iotdataready.movistar.es

One can find more information on how to configure this particular device on [the product's wiki][doc-main].
Specific sites exist for [setting up the UART connection][doc-uart] and for configuring an
[upstream connection][doc-upstream].

1nce also provides a wealth of information, including [how to configure its APN][doc-apn].

### Flashing the FW
One can flash new firmware versions to the devices in order to add new features or fix outstanding errors. The
information on how to do that can be found [here][dragino-uart]. The bottom line is one has to pull the
appropriate firmware and flash it with the [STM32 Cube Programmer][stm32-programmer].

The firmware's flashed on an address specified by the user. The documentation is not clear about what these should
be, so bear in mind that:

- When flashing the **bootloader** the starting address should be `0x8000000`.
- When flashing the **firmware** the starting address should be `0x8007800`.

If anything's flashed in the bootloader's address space, the bootloader should be reflashed before attempting
to flash the firmware itself.

The applicable bootlader and firmware for our SE0Xs are located on the `firmware/` directory.

## Deploying the Lambda
Every action one can take regarding the contents of this repository in terms of deployment has been automated
by means of a `Makefile`. You can just invoke `make(1)` and go on from there!

<!-- REFs -->
[doc-main]: https://wiki.dragino.com/xwiki/bin/view/Main/User%20Manual%20for%20LoRaWAN%20End%20Nodes/SE0X-NBNS--NB-IoT_Soil%20Moisture_%26_EC_Sensor_Transmitter_User_Manual/#H1.1WhatisSE0X-NB2FNSNB-IoTSoilMoisture26ECSensor
[doc-uart]: https://wiki.dragino.com/xwiki/bin/view/Main/UART_Access_for_NB_ST_BC660K-GL/#H4.2UpdateFirmware28Assumethedevicealreadyhaveabootloader29
[doc-calibration]: https://www.dragino.com/downloads/downloads/LoRa_End_Node/LSE01/Calibrate_to_other_Soil_20230522.pdf
[doc-upstream]: https://wiki.dragino.com/xwiki/bin/view/Main/General%20Configure%20to%20Connect%20to%20IoT%20server%20for%20-NB%20%26%20-NS%20NB-IoT%20models/#H2.AttachNetwork
[doc-apn]: https://help.1nce.com/dev-hub/docs/data-services-apn
[doc-once-api]: https://help.1nce.com/dev-hub/reference/api-welcome
[doc-once-api-auth]: https://help.1nce.com/dev-hub/reference/postaccesstokenpost
[doc-once-webhook-api]: https://help.1nce.com/dev-hub/docs/cloud-integrator-webhook-configuration
[tsdb-create-hypertable]: https://docs.tigerdata.com/api/latest/hypertable/create_table/
[lego-renewal]: https://go-acme.github.io/lego/usage/cli/renew-a-certificate/index.html
[tsdb-write-data]: https://docs.tigerdata.com/use-timescale/latest/write-data/insert/
[doc-go-pq]: https://pkg.go.dev/github.com/lib/pq
[doc-go-sql]: https://pkg.go.dev/database/sql
[stm32-programmer]: https://www.st.com/en/development-tools/stm32cubeprog.html
