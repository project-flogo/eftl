FROM debian
RUN apt-get update && apt-get install -y default-jre wget unzip
RUN wget "https://edownloads.tibco.com/Installers/tap/FTL-CE/6.0.1/TIB_ftl-ce_6.0.1_linux_x86_64.zip?SJCDPTPG=1554670084_c00e5d2aa2eeaca4abd002ac4ab7f5cb&ext=.zip" \
  -O TIB_ftl-ce_6.0.1_linux_x86_64.zip
RUN wget "https://edownloads.tibco.com/Installers/tap/EFTL-CE/6.0.1/TIB_eftl-ce_6.0.1_linux_x86_64.zip?SJCDPTPG=1554670086_3d0242df9809cf1cef9d0cb1d4125561&ext=.zip" \
  -O TIB_eftl-ce_6.0.1_linux_x86_64.zip
RUN unzip TIB_ftl-ce_6.0.1_linux_x86_64.zip
RUN unzip TIB_eftl-ce_6.0.1_linux_x86_64.zip
RUN dpkg -i TIB_FTL_6.0.1/deb/*.deb
RUN dpkg -i TIB_EFTL_6.0.1/deb/*.deb
RUN rm -rf TIB_FTL_6.0.1/ TIB_EFTL_6.0.1/ TIB_ftl-ce_6.0.1_linux_x86_64.zip TIB_eftl-ce_6.0.1_linux_x86_64.zip
RUN echo "/opt/tibco/ftl/current-version/bin/tibftlserver --config /opt/tibco/eftl/6.0/samples/tibftlserver_eftl.yaml --name SRV1&" > run.sh
RUN echo "sleep 30" >> run.sh
RUN echo "/opt/tibco/eftl/6.0/ftl/bin/tibftladmin --ftlserver http://localhost:8585 --updaterealm /opt/tibco/eftl/6.0/samples/tibrealm.json" >> run.sh
RUN echo "/opt/tibco/eftl/6.0/ftl/bin/tibftlserver --config /opt/tibco/eftl/6.0/samples/tibftlserver_eftl.yaml --name EFTL" >> run.sh
CMD ["/bin/bash", "run.sh"]
