//
//  FormTextInputCell.swift
//  Permission Hub
//
//  Created by Corné on 13/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import UIKit

final class FormTextInputCell: UITableViewCell {

    // MARK: - Private properties

    private lazy var textField: UITextField = {

        let textField = UITextField()
        textField.translatesAutoresizingMaskIntoConstraints = false
        textField.font = PHFonts.regular(ofSize: 14)
        textField.textColor = PHColors.greyishBrown
        textField.tintColor = PHColors.greyishBrown
        textField.clearButtonMode = .whileEditing
        textField.backgroundColor = PHColors.formGray
        textField.isUserInteractionEnabled = false

        textField.addTarget(
            self,
            action: #selector(textFieldDidChange),
            for: .editingChanged)

        return textField
    }()

    private var callback: ((String) -> Void)?

    // MARK: - Initialization

    override init(style: UITableViewCell.CellStyle, reuseIdentifier: String?) {
        super.init(style: style, reuseIdentifier: reuseIdentifier)

        configure()
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    // MARK: - Configuration

    func configure(withPlaceholder
        placeholder: String, andText
        text: String? = nil, andKeyboardType
        keyboardType: UIKeyboardType,
        callback: @escaping (String) -> Void) {

        textField.attributedPlaceholder = NSAttributedString(
            string: placeholder,
            attributes: [.foregroundColor: PHColors.greyishBrown.withAlphaComponent(0.5)])
        textField.text = text
        textField.keyboardType = keyboardType

        self.callback = callback
    }

    private func configure() {

        addSubview(textField)

        textField.heightAnchor.constraint(equalToConstant: 44).isActive = true
        textField.topAnchor.constraint(equalTo: topAnchor, constant: 5).isActive = true
        textField.leftAnchor.constraint(equalTo: leftAnchor, constant: 20).isActive = true
        textField.rightAnchor.constraint(equalTo: rightAnchor, constant: -20).isActive = true
        textField.bottomAnchor.constraint(equalTo: bottomAnchor, constant: -5).isActive = true
    }

    // MARK: - Selectors

    @objc private func textFieldDidChange(_ sender: UITextField) {

        guard let text = sender.text else {
            return
        }

        callback?(text)
    }
}
